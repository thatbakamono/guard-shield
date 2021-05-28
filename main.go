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
	"%{color}%{time:15:04:05} %{shortfunc} %{level:.4s} %{color:reset}%{message}",
)

func main() {
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)

	logging.SetBackend(backendFormatter)

	err := redirect(sourceAddress, targetAddress)

	if err != nil {
		panic(err)
	}
}

func redirect(sourceAddress string, targetAddress string) error {
	listener, err := net.Listen("tcp", sourceAddress)

	if err != nil {
		return err
	}

	defer listener.Close()

	log.Infof("Listening on %s", sourceAddress)

	for {
		clientConnection, err := listener.Accept()

		if err != nil {
			log.Error("Failed to accept connection")
			log.Error(err)

			continue
		}

		log.Info("Accepted connection")

		go handleClient(clientConnection, targetAddress)
	}
}

func handleClient(clientConnection net.Conn, targetAddress string) {
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
