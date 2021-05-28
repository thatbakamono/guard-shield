package main

import (
	"fmt"
	"io"
	"net"
)

const sourceAddress = "127.0.0.1:25565"
const targetAddress = "127.0.0.1:25566"

func main() {
	fmt.Println("Hello guard-shield")

	listener, err := net.Listen("tcp", sourceAddress)

	if err != nil {
		panic(err)
	}

	defer listener.Close()

	fmt.Println("Listening")

	for {
		connection, err := listener.Accept()

		if err != nil {
			panic(err)
		}

		fmt.Println("Accepted connection")

		go handleClient(connection)
	}
}

func handleClient(clientConnection net.Conn) {
	defer clientConnection.Close()

	serverConnection, err := net.Dial("tcp", targetAddress)

	if err != nil {
		panic(err)
	}

	defer serverConnection.Close()

	go io.Copy(clientConnection, serverConnection)
	io.Copy(serverConnection, clientConnection)
}
