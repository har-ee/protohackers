package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

func insert(timestamp int32, price int32, prices map[int32]int32) {
	prices[timestamp] = price
}

func query(mintime int32, maxtime int32, prices map[int32]int32) int64 {
	sum := int64(0)
	n := int64(0)

	for time, price := range prices {
		if time >= mintime && time <= maxtime {
			sum += int64(price)
			n++
		}
	}

	if n == 0 {
		return 0
	}

	return int64(sum / n)
}

func handler(conn net.Conn, clientid int) {
	defer conn.Close()

	prices := make(map[int32]int32)

	for {
		buf := make([]byte, 9)
		_, err := io.ReadFull(conn, buf)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("[%d] Error readin: %#v\n", clientid, err)
			}
			return
		}

		command := rune(buf[0])
		i1 := int32(binary.BigEndian.Uint32(buf[1:5]))
		i2 := int32(binary.BigEndian.Uint32(buf[5:9]))

		if command == 'I' {
			insert(i1, i2, prices)
		} else {
			mean := query(i1, i2, prices)

			responseBytes := make([]byte, 4)
			binary.BigEndian.PutUint32(responseBytes, uint32(mean))

			_, err := conn.Write(responseBytes)
			if err != nil {
				fmt.Printf("Error writing: %#v\n", err)
				return
			}
		}

		fmt.Printf("[%d] Parsed: %v, %d, %d\n", clientid, string(command), i1, i2)
	}
}

func main() {
	fmt.Println("Passing over to handler")

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("Error listening: %#v\n", err)
	}

	for i := 0; true; i++ {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Error accepting: %#v\n", err)
		}

		go handler(conn, i)
	}
}
