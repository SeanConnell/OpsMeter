package main

import (
	"bytes"
	"github.com/hoisie/web"
	"github.com/tarm/goserial"
	"io"
	"log"
	"fmt"
	"encoding/hex"
)

/*
Runs a small webserver that communicates with an embedded lighting controller.
*/

//TODO: make this format a packet with data in the format COMMAND\r[DATA BYTES]
func send(command string, serial io.ReadWriteCloser) {
	serial.Write([]byte(command + "\r"))
}

func recieve(size int, serial io.ReadWriteCloser) ([]byte , error){
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
	if err != nil{
		ctx.WriteString(err.Error())
	} else {
		out := hex.EncodeToString(msg)
		ctx.WriteString(out)
	}
}

//Used to handle getting input data to set 
//TODO implement
func input_data(ctx *web.Context, serial io.ReadWriteCloser) {
	retval := "SETOUPUTDATA: "
	for k, v := range ctx.Params {
		retval += k + "->" + v + ","
	}
	retval += "\n"
	var buf bytes.Buffer
	buf.WriteString(retval)
	send("SETOUTPUTDATA", serial)
	io.Copy(ctx, &buf)
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
		"<li>POST to /SETOUTPUT with 32 3-tuples of 8 bit [R,G,B] values in a format like:<br>" +
		"    0R=138&0G=12&0B=241&1R=23& etc ...<br>" +
		"    Any values not set will default to 0,0,0 eg black.</li>" +
		"<li>POST to /RESET will send a reset command to the light controller and reset state.</li>" +
		"<li>GET to /STATE will read device state and display it here.</li>" +
		"<li>GET to /* (anything but /STATE) will display this command list.</li>" +
		"</ul>" +
		"<br>" +
		"<i>Created and maintained by connells</i>"
	return commands
}

//Setup handlers for different addresses and HTTP methods here
func main() {
	//setup serial for use
	s, err := initialize_serial("/dev/ttyACM0", 115200)
	if err != nil {
		log.Fatal("Couldn't open comm port. ", err)
	}
	defer s.Close()

	//close over lower level function to prevent reinitialization of serial port
	send_input_data := func(ctx *web.Context) { input_data(ctx, s) }
	send_reset := func(ctx *web.Context) { reset(ctx, s) }
	get_output_state := func(ctx *web.Context) {output_state(ctx, s) }

	//setup handlers for various commands/addresses
	web.Get("/STATE", get_output_state)         //Get state from device and return it
	web.Get("/.*", commands_index)          //Get help
	web.Put("/STATE", send_input_data) //Set the output of the device
	web.Put("/RESET", send_reset)          //Reset device
	web.Run("0.0.0.0:8089")
}
