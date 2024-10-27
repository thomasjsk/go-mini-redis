package main

import (
	"errors"
	"fmt"
	"github.com/thomasjsk/go-mini-redis/command"
	"github.com/thomasjsk/go-mini-redis/frame"
	"io"
	"net"
	"os"
)

func main() {
	server := createServer()
	server.Start()
}

type Server struct {
	listener net.Listener
	quitch   chan struct{}
}

func createServer() *Server {
	return &Server{
		quitch: make(chan struct{}),
	}
}

func (server *Server) Start() {
	listener, err := net.Listen("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	fmt.Println("Redis server running on port 6379")
	defer listener.Close()

	server.listener = listener
	go server.acceptConnections()
	<-server.quitch
}

func (server *Server) acceptConnections() {
	for {
		conn, err := server.listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		fmt.Println("New connection accepted")
		go server.readCommands(conn)
	}
}
func (server *Server) readCommands(conn net.Conn) {
	defer conn.Close()

	for {
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if errors.Is(err, io.EOF) {
			fmt.Println("Client closed the connections:", conn.RemoteAddr())
			break
		} else if err != nil {
			fmt.Println("Error while reading the message")
		}

		rawRequestFrame := frame.RawFrame(buf)
		fmt.Println("rawRequestFrame", string(buf))

		requestFrame, err := rawRequestFrame.Decode(new(int))
		if err != nil {
			fmt.Println("Error while decoding the cursor", err)
		}
		fmt.Println("requestFrame", requestFrame)

		responseFrame := HandleCommand(command.Command{}.FromFrame(&requestFrame))

		fmt.Println("REQUEST :", requestFrame)
		fmt.Println("RESPONSE:", responseFrame)

		_, err = conn.Write(responseFrame.Encode())
		if err != nil {
			fmt.Println("Error while writing the response", err)
		}

		fmt.Println("____________")

	}
}
