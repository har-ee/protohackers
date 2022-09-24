package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net"
)

type Request struct {
	Method *string  `json:"method"`
	Number *float64 `json:"number"`
}

type Response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
	Error  string `json:"error"`
}

func parse(input []byte) (Request, error) {
	var r Request

	err := json.Unmarshal(input, &r)

	if err != nil {
		fmt.Printf("Error unmarshaling input: %#v\n", err)
		return r, err
	} else if r.Method == nil {
		return r, errors.New("Method not given")
	} else if r.Number == nil {
		return r, errors.New("Number not given")
	}

	return r, nil
}

func handleCommand(request Request) Response {
	var response Response
	if *request.Method == "isPrime" {
		response = Response{
			Method: "isPrime",
			Prime:  big.NewInt(int64(*request.Number)).ProbablyPrime(0),
		}
	} else {
		response = Response{
			Error: "Unsupported Method",
		}
	}
	return response
}

func sendResponse(conn net.Conn, response Response) {
	rawResponse, err := json.Marshal(response)
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Printf("Responding with: %s\n", string(rawResponse))

	_, err = conn.Write(append(rawResponse, '\n'))
	if err != nil {
		fmt.Printf("Error writing: %#v\n", err)
	}
}

func handler(conn net.Conn) {
	defer conn.Close()

	var reader = bufio.NewReader(conn)

	for {
		bytes, err := reader.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Error readin: %#v\n", err)
			}
			return
		}

		fmt.Printf("Message received: %s\n", string(bytes))

		request, err := parse(bytes)
		if err != nil {
			response := Response{
				Error: "Invalid JSON",
			}
			sendResponse(conn, response)
			return
		}

		response := handleCommand(request)
		sendResponse(conn, response)
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
