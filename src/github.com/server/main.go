package main

import (
	"fmt"
	"net"
)

func main() {
	var _ = fmt.Printf

	ln, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(connection net.Conn) {
	str := make([]byte, 1024)

	for {
		if _, err := connection.Read(str); err != nil {
			panic(err)
		}

		if string(str[:4]) == "exit" {
			if err := connection.Close(); err != nil {
				panic(err)
			}
		}
	}
}