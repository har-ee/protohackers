package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"regexp"
)

func broadcast(broadcasts chan Broadcast, username string, message string) {
	fmt.Printf("[%s] Sending broadcast, %s\n", username, message)
	broadcasts <- Broadcast{
		username: username,
		message:  message,
	}
}

func disconnect(disconnections chan Disconnection, username string) {
	e := Disconnection{username}
	disconnections <- e
}

func (comm *Communication) frontend(conn net.Conn) {
	defer conn.Close()
	var reader = bufio.NewReader(conn)

	fmt.Printf("New client connected\n")

	username, err := getUsername(conn)
	if err != nil {
		return
	}

	messages := make(chan string)
	register(username, comm.registrations, messages)

	go messageListener(messages, conn)

	// Main loop
	for {
		bytes, err := reader.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Error readin: %#v\n", err)
			}

			disconnect(comm.disconnections, username)
			return
		}

		message := string(bytes)

		fmt.Printf("Message received: %s\n", message)

		broadcast(comm.broadcasts, username, message)
	}
}

func getUsername(conn net.Conn) (string, error) {
	var reader = bufio.NewReader(conn)
	_, err := conn.Write([]byte("Welcome to hchat! Please enter a name\n"))
	if err != nil {
		fmt.Printf("Error writing: %#v\n", err)
		return "", err
	}

	nameBytes, err := reader.ReadBytes('\n')
	if err != nil {
		if err != io.EOF {
			fmt.Printf("Error readin: %#v\n", err)
		}
		return "", err
	}
	name := string(nameBytes)
	name = name[:len(name)-1] // Strip trailing /n

	// Check name
	m, _ := regexp.MatchString("^[a-zA-Z0-9]{1,16}$", name)
	if !m {
		_, err = conn.Write([]byte("Error: invalid name!"))
		if err != nil {
			fmt.Printf("Error writing: %#v\n", err)
		}
		return "", errors.New("Invalid name")
	}
	return name, nil
}

func messageListener(messages chan string, conn net.Conn) {
	for {
		message := <-messages
		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Printf("Error writing: %#v\n", err)
			return
		}
	}

}

func register(username string, registrations chan Registration, messages chan string) {
	registrations <- Registration{
		username: username,
		messages: messages,
	}
}
