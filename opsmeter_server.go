package main

import (
    "github.com/hoisie/web"
    "github.com/tarm/goserial"
    "bytes"
    "io"
    "log"
)

/*
Runs a small webserver that communicates with an embedded lighting controller.
TODO:Dont' do so many syscalls for files everywhere, do it once and hold state
*/

//Tell the device to reset itself
//TODO implement
func reset(ctx *web.Context) string { 
    var message = 
    "TODO: Reset device"
    return message 
}

//Get state from device and display it
//TODO implement
func output_state(ctx *web.Context) string { 
    var state = 
    "TODO: communicate with device"
    return state 
}   

//Used to handle getting input data to set 
//TODO implement
func input_data(ctx *web.Context) { 
    retval := "Input: "
    for k,v := range ctx.Params {
        retval = retval + k + "->" + v + ","
    }
    var buf bytes.Buffer
    buf.WriteString(retval)
    c := &serial.Config{Name:"/dev/ttyACM0", Baud:115200}
    s, err := serial.OpenPort(c)
    defer s.Close()
    if err!= nil {
        log.Fatal(err)
    }
    s.Write([]byte("test"))

    io.Copy(ctx, &buf)
}   

//Display a help message to a user about how to interact with this service
func commands_index(ctx *web.Context) string { 
    var commands = 
    "Commands index for <b>OPS LIGHTS</b>" +
    "<ul>" + 
    "<li>POST to /SETOUTPUT with 32 3-tuples of 8 bit [R,G,B] values in a format like:<br>"+
    "    0R=138&0G=12&0B=241&1R=23& etc ...<br>" +
    "    Any values not set will remain in their previous state.</li>" +
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
    web.Get("/STATE", output_state) //Get state from device and return it
    web.Get("/.*", commands_index) //Get help
    web.Post("/SETOUTPUT", input_data) //Set the output of the device
    web.Post("/RESET", reset) //Reset device
    web.Run("0.0.0.0:9999")
}
