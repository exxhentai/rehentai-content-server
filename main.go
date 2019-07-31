package main

import (
	"comm"
	"flag"
)

func main() {
	port := flag.Int("port", 1234, "")

	flag.Parse()

	server := comm.NewServer()
	server.Start(*port)
}
