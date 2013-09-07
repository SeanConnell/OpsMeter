package main

import (
    "github.com/hoisie/web"
    "fmt"
)

func hello(ctx *web.Context, val string) string { 
    for k,v := range ctx.Params {
        println(k, v)
    }
    return "Lights server!"
}   

func data(ctx *web.Context, val string) string { 
    retval := "Input: "
    for k,v := range ctx.Params {
        retval = retval + k + "->" + v + ","
    }
    fmt.Println(val)
    return val
}   

func main() {
    web.Get("/(.*)", hello)
    web.Post("/(.*)", data)
    web.Run("0.0.0.0:9999")
}
