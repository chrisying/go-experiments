package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error on Listen: ", err)
		os.Exit(-1)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error on Accept: ", err)
			os.Exit(-1)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	var buffer [1024]byte
	for {
		n, err := conn.Read(buffer[:])
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed")
				conn.Close()
			} else {
				fmt.Println("Error on Read: ", err)
			}
		}
		msg := string(buffer[:n])
		fmt.Println("Client sent: ", msg)
	}
}
