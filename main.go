package main

import (
	"github.com/op/go-logging"
	"io"
	"net"
	"os"
)

const sourceAddress = "127.0.0.1:25565"
const targetAddress = "127.0.0.1:25566"

var log = logging.MustGetLogger("guard-shield")
var format = logging.MustStringFormatter(
	"%{color}%{time:15:04:05} %{shortfunc} %{level:.4s} %{color:reset}%{message}"
)

func main() {
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)

	logging.SetBackend(backendFormatter)

	listener, err := net.Listen("tcp", sourceAddress)

	if err != nil {
		panic(err)
	}

	defer listener.Close()

	log.Info("Listening")

	for {
		connection, err := listener.Accept()

		if err != nil {
			panic(err)
		}

		log.Info("Accepted connection")

		go handleClient(connection)
	}
}

func handleClient(clientConnection net.Conn) {
	defer clientConnection.Close()
	defer log.Info("Closed connection")

	serverConnection, err := net.Dial("tcp", targetAddress)

	if err != nil {
		panic(err)
	}

	defer serverConnection.Close()

	go io.Copy(clientConnection, serverConnection)
	io.Copy(serverConnection, clientConnection)
}
