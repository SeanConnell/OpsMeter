package main

import (
    "github.com/hoisie/web"
    "bytes"
    "io"
)

func commands_index(ctx *web.Context, val string) string { 
    var commands = 
    "Commands index for OPS LIGHTS:<br>" +
    "POST to /SETOUTPUT with 32 values of [R,G,B] values in a format like:<br>"+
    "0R=138&0G=12&0B=241&1R=23& etc ...<br>" +
    "Any values not set will remain in their previous state.<br>" +
    "Created and maintained by connells<br>"
    return commands 
}   

func input_data(ctx *web.Context) { 
    retval := "Input: "
    for k,v := range ctx.Params {
        retval = retval + k + "->" + v + ","
    }
    var buf bytes.Buffer
    buf.WriteString(retval)
    io.Copy(ctx, &buf)
}   

func main() {
    web.Get("/(.*)", commands_index)
    web.Post("/SETOUTPUT", input_data)
    web.Run("0.0.0.0:9999")
}
