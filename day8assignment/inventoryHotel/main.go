package main

import (
	"flag"
	"fmt"
)


func main() {
	//passing port and source name to flags from terminal
	var port = flag.Int("p", 8081, "please specify the port to run")
	var source = flag.String("b", "Source-A", "Please specify source")
	flag.Parse()
	fmt.Println("Starting server as ", *source, " on port ", *port)
	
	routing(*port)
}
