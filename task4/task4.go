package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

var (
	store = make(map[string]string)
	mu    sync.Mutex
)

func handler(packet string, conn net.PacketConn, addr net.Addr) {
	var key string
	var value string

	s := strings.SplitN(packet, "=", 2)

	key = s[0]

	if len(s) > 1 {
		value = s[1]
		insert(key, value)
	} else {
		value = retrieve(key)
		response := fmt.Sprintf("%s=%s", key, value)
		conn.WriteTo([]byte(response), addr)
	}

}

func insert(key string, value string) {
	if key == "version" {
		return
	}

	mu.Lock()

	store[key] = value

	mu.Unlock()

	fmt.Printf("Set: %s=%s\n", key, value)
}

func retrieve(key string) string {
	value := store[key]

	fmt.Printf("Retrieved: %s=%s\n", key, value)

	return value
}

func main() {

	store["version"] = "My key store!"

	ln, err := net.ListenPacket("udp", ":8080")
	if err != nil {
		fmt.Printf("Error listening: %#v\n", err)
	}

	// Accept new connections
	for {
		buf := make([]byte, 1000)
		n, addr, err := ln.ReadFrom(buf)
		if err != nil {
			fmt.Printf("Error accepting: %#v\n", err)
		}
		go handler(string(buf[:n]), ln, addr)
	}
}
