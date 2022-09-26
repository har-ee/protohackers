package main

import (
	"fmt"
	"strings"
)

func (c *Communication) backend() {
	userChans := make(map[string]chan string)

	for {
		select {
		case broadcast := <-c.broadcasts:
			doBroadcast(broadcast, userChans)
		case registration := <-c.registrations:
			doRegistrations(registration, userChans)
		case disconnection := <-c.disconnections:
			doDisconnection(disconnection, userChans)
		}
	}
}

func doBroadcast(broadcast Broadcast, chans map[string]chan string) {
	fmt.Printf("Broadcasting [%s] %s\n", broadcast.username, broadcast.message)
	for user, channel := range chans {
		if user != broadcast.username {
			channel <- fmt.Sprintf("[%s] %s", broadcast.username, broadcast.message)
		}
	}
}

func doDisconnection(disconnection Disconnection, chans map[string]chan string) {
	fmt.Printf("User disconnected %s\n", disconnection.username)
	delete(chans, disconnection.username)
	for _, channel := range chans {
		channel <- fmt.Sprintf("* %s has left the room\n", disconnection.username)
	}
}

func doRegistrations(registration Registration, chans map[string]chan string) {
	fmt.Printf("Registering user %s\n", registration.username)
	var users []string

	for user, channel := range chans {
		channel <- fmt.Sprintf("* %s has entered the room\n", registration.username)
		users = append(users, user)
	}

	chans[registration.username] = registration.messages
	registration.messages <- fmt.Sprintf("* The room contains: %s\n", strings.Join(users, ", "))
}
