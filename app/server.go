package main

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/thomasjsk/go-mini-redis/command"
	"github.com/thomasjsk/go-mini-redis/frame"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	server := createServer()
	server.Start()
}

type Server struct {
	listener    net.Listener
	connections map[uuid.UUID]*Connection
	config      ServerConfig
	quitCh      chan struct{}
}

func createServer() *Server {
	server := &Server{
		connections: make(map[uuid.UUID]*Connection),
		quitCh:      make(chan struct{}),
		config: ServerConfig{
			port:       6379,
			replOffset: 0,
			RDB:        ServerConfigRDB{},
		},
	}

	return server
}

func (server *Server) Start() {
	fmt.Println("-------------------------------START SERVER-------------------------------")
	server.Configure()
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", server.config.port))
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	fmt.Printf("Redis server running on port %d\n", server.config.port)
	defer listener.Close()

	server.listener = listener
	go server.acceptConnections()
	go server.replicateFromMaster(&server.config)
	<-server.quitCh
}

func (server *Server) replicateToSlaves(cmd command.Command, f frame.Frame) {
	fmt.Println("replicateToSlaves")
	replicate := false

	switch cmd.Name {
	case "set":
		replicate = true
	}

	if !replicate {
		fmt.Println("Command not replicable: ", cmd.Name)
		return
	}

	for _, conn := range server.connections {
		fmt.Println("Looping over connections", conn)
		if conn.role == Slave {
			fmt.Println("Replicating to slaves")
			_, err := conn.Write(f)
			if err != nil {
				fmt.Println("Failed to write to slaves")
			}
		}
	}
}

func (server *Server) acceptConnections() {
	for {
		conn, err := server.listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		fmt.Println("New connection accepted")
		connectionId := uuid.New()

		connection := &Connection{
			id:         connectionId,
			role:       Master,
			connection: conn,
			onCloseCh:  make(chan int),
		}

		server.connections[connectionId] = connection
		go server.listen(connection)
	}
}
func (server *Server) listen(connection *Connection) {
	defer connection.Close()

	isMaster := server.config.replicaOf.port == 0
	var mode string
	if isMaster {
		mode = "MASTER"
	} else {
		mode = "SLAVE"
	}

	fmt.Printf("[%s] listen...\n", mode)

	go func() {
		<-connection.onCloseCh
		delete(server.connections, connection.id)
		close(connection.onCloseCh)
	}()

	for {
		buf := make([]byte, 1024)
		_, err := connection.Read(buf)
		if errors.Is(err, io.EOF) {
			fmt.Println("Client closed the connection")
			connection.onCloseCh <- 0
			break
		} else if err != nil {
			fmt.Println("Error while reading the message: ", err)
			connection.onCloseCh <- 1
			break
		}

		rawRequest := frame.RawFrame(buf)
		req := strings.ReplaceAll(string(buf), "\n", "\\n")
		req = strings.ReplaceAll(req, "\r", "\\r")
		fmt.Printf("[%s] received %s\n", mode, req)

		request, err := rawRequest.Decode(new(int))
		if err != nil {
			fmt.Println("Error while decoding the cursor", err)
		}

		cmd := command.Command{}.FromFrame(&request)
		err = HandleCommand(cmd, connection, server, false)

		if err != nil {
			fmt.Println("Error while handling the command", err)
		}

		server.replicateToSlaves(cmd, request)
	}
}

func (server *Server) replicateFromMaster(config *ServerConfig) {
	if len(config.replicaOf.host) == 0 {
		return
	}

	conn := config.replicaOf.connection

	steps := []frame.Frame{
		frame.CreateArray([]frame.Frame{
			frame.CreateBulkString("PING"),
		}),
		frame.CreateArray([]frame.Frame{
			frame.CreateBulkString("REPLCONF"),
			frame.CreateBulkString("listening-port"),
			frame.CreateBulkString(strconv.Itoa(server.config.port)),
		}),
		frame.CreateArray([]frame.Frame{
			frame.CreateBulkString("REPLCONF"),
			frame.CreateBulkString("capa"),
			frame.CreateBulkString("eof"),
			frame.CreateBulkString("capa"),
			frame.CreateBulkString("psync2"),
		}),
		frame.CreateArray([]frame.Frame{
			frame.CreateBulkString("PSYNC"),
			frame.CreateBulkString("?"),
			frame.CreateBulkString("-1"),
		}),
	}

	for _, step := range steps {
		_, err := conn.Write(step)
		buf := make([]byte, 1024)
		conn.Read(buf)
		req := strings.ReplaceAll(string(buf), "\n", "\\n")
		req = strings.ReplaceAll(req, "\r", "\\r")
		fmt.Printf("Received %s\n", req)
		if err != nil {
			fmt.Println("Error while writing to connection: ", err)
		}
	}

	fmt.Println("Replicating...")
	for {
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if errors.Is(err, io.EOF) {
			fmt.Println("Client closed the connection")
			conn.onCloseCh <- 0
			break
		} else if err != nil {
			fmt.Println("Error while reading the message: ", err)
			conn.onCloseCh <- 1
			break
		}

		rawRequest := frame.RawFrame(buf)
		req := strings.ReplaceAll(string(buf), "\n", "\\n")
		req = strings.ReplaceAll(req, "\r", "\\r")
		fmt.Printf("[REPLICATION] received %s\n", req)

		request, err := rawRequest.Decode(new(int))
		if err != nil {
			fmt.Println("Error while decoding the cursor", err)
		}

		err = HandleCommand(command.Command{}.FromFrame(&request), &conn, server, true)

		if err != nil {
			fmt.Println("Error while handling the command", err)
		}
	}
}
