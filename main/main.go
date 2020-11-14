package main

import (
	"flag"
)

var(
	httpListenAddr = flag.String("httpListenAddr", ":8481", "Address to listen for http connections")
)

func main(){
}
