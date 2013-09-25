package main

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/hoisie/web"
	"github.com/tarm/goserial"
	"io"
	"log"
	"strconv"
)

/*
Runs a small webserver that communicates with an embedded lighting controller.
*/

func send(data string, serial io.ReadWriteCloser) {
	serial.Write([]byte(data + "\r"))
}

func recieve(size int, serial io.ReadWriteCloser) ([]byte, error) {
	b := make([]byte, size)
	n, err := serial.Read(b)
	//don't handle reading fewer bytes than expected for now
	if n < size {
		return nil, fmt.Errorf("Expected %v bytes and read %v instead", size, n)
	}
	if err != nil {
		return nil, err
	}
	return b, nil
}

//Tell the device to reset itself
//TODO implement error checking
func reset(ctx *web.Context, serial io.ReadWriteCloser) {
	send("RESET", serial)
	retval := "RESET sent\n"
	var buf bytes.Buffer
	buf.WriteString(retval)
	io.Copy(ctx, &buf)
}

//Get state from device and display it
func output_state(ctx *web.Context, serial io.ReadWriteCloser) {
	send("GETSTATE", serial)
	msg, err := recieve(96, serial)
	if err != nil {
		ctx.WriteString(err.Error())
	} else {
		out := hex.EncodeToString(msg)
		ctx.WriteString(out)
	}
}

//Used to handle getting input data to set 
func input_data(ctx *web.Context, serial io.ReadWriteCloser) {
	retval := "SETSTATE: "
	output_state := make([]byte, 96)
	for k, v := range ctx.Params {
		retval += k + "->" + v + ","
		color_loc, err := parse_color_location(k)
		if err != nil {
			ctx.WriteString(err.Error())
			return
		}
		if len(v) != 1 {
			ctx.WriteString("Value must be 1 byte long, eg between 0-255")
		}
		output_state[color_loc] = ([]byte(v))[0] //sketchy cast, fix with validation later
	}
	retval += string(output_state)
	retval += "\n"
	var buf bytes.Buffer
	buf.WriteString(retval)
	send("SETSTATE"+string(output_state), serial)
	io.Copy(ctx, &buf)
}

func parse_color_location(color_index string) (uint64, error) {
	if len(color_index) != 3 {
		return 0, errors.New("Improper color index length of not 3")
	}
	color_pos, err := get_pos(color_index[len(color_index)-1:])
	if err != nil {
		return 0, err
	}
	//convert string to unsigned base 10 int from 8 bits
	led_pos, err := strconv.ParseUint((color_index[:len(color_index)-1]), 10, 8)
	if err != nil {
		return 0, err
	}
	if led_pos >= 32 {
		return 0, errors.New("The led strip is only 32 leds long, zero indexed")
	}
	//* 3 because each LED is three wide
	return (led_pos * 3) + uint64(color_pos), nil
}

func get_pos(name string) (int, error) {
	switch name {
	case "r", "R":
		return 0, nil
	case "g", "G":
		return 1, nil
	case "b", "B":
		return 2, nil
	default:
		return 0, errors.New("color position must end with r,g,b (case insensitive)")
	}
	return 0, errors.New("This should never happen")
}

//Abstract/encapsulate library for serial usage
func initialize_serial(name string, baud int) (io.ReadWriteCloser, error) {
	c := &serial.Config{Name: name, Baud: baud}
	s, err := serial.OpenPort(c)
	if err != nil {
		return nil, err
	}
	return s, nil
}

//Display a help message to a user about how to interact with this service
func commands_index(ctx *web.Context) string {
	var commands = "Commands index for <b>OPS LIGHTS</b>" +
		"<ul>" +
		"<li>GET to /STATE will read device state and display it here.</li>" +
		"<li>POST to /STATE with 32 3-tuples of 8 bit [R,G,B] values in a format like:<br>" +
		"    0R=138&0G=12&0B=241&1R=23& etc ...<br>" +
		"    The whole array must be supplied.</li>" +
		"<li>POST to /RESET will send a reset command to the light controller and reset state.</li>" +
		"<li>GET to /* (anything but /STATE) will display this command list.</li>" +
		"</ul>" +
		"<br>" +
		"<i>Created and maintained by connells</i>"
	return commands
}

//Setup handlers for different addresses and HTTP methods here
func main() {
	//setup serial for use
	s, err := initialize_serial("/dev/ttyACM1", 115200)
	if err != nil {
		log.Fatal("Couldn't open comm port. ", err)
	}
	defer s.Close()

	//close over lower level function to prevent reinitialization of serial port
	send_input_data := func(ctx *web.Context) { input_data(ctx, s) }
	send_reset := func(ctx *web.Context) { reset(ctx, s) }
	get_output_state := func(ctx *web.Context) { output_state(ctx, s) }

	//setup handlers for various commands/addresses
	web.Get("/STATE", get_output_state) //Get state from device and return it
	web.Get("/.*", commands_index)      //Get help
	web.Put("/STATE", send_input_data)  //Set the output of the device
	web.Put("/RESET", send_reset)       //Reset device
	web.Run("0.0.0.0:8089")
}
