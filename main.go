package main

import (
	"comm"
)

func main() {
	server := comm.NewServer()
	server.Start(1234)
}
