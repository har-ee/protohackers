package main

import (
	"fmt"
	"net"
)

type Broadcast struct {
	username string
	message  string
}

type Communication struct {
	broadcasts     chan Broadcast
	registrations  chan Registration
	disconnections chan Disconnection
}

type Disconnection struct {
	username string
}

type Registration struct {
	username string
	messages chan string
}

func main() {
	c := &Communication{
		broadcasts:     make(chan Broadcast),
		disconnections: make(chan Disconnection),
		registrations:  make(chan Registration),
	}

	// Start chatroom backend
	go c.backend()

	fmt.Println("Passing over to handler")

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("Error listening: %#v\n", err)
	}

	// Accept new connections
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Error accepting: %#v\n", err)
		}
		go c.frontend(conn)
	}
}
