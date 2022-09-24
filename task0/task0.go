package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
)

func handler(conn net.Conn) {
	defer conn.Close()

	var buf bytes.Buffer

	len, err := io.Copy(&buf, conn)
	if err != nil {
		fmt.Printf("Error readin: %#v\n", err)
		return
	}

	fmt.Printf("Message received: %s\n", string(buf.Bytes()[:len]))

	_, err = conn.Write(buf.Bytes())
	if err != nil {
		fmt.Printf("Error writing: %#v\n", err)
		return
	}
}

func main() {
	fmt.Println("Passing over to handler")

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("Error listening: %#v\n", err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Error accepting: %#v\n", err)
		}
		go handler(conn)
	}
}
