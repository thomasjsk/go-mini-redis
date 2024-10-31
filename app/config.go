package main

import (
	"fmt"
	"github.com/google/uuid"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

type ServerConfigRDB struct {
	Dir        string
	DBFileName string
}

type ServerConfigReplicaOf struct {
	host       string
	port       int
	connection Connection
}

type ServerConfig struct {
	mu         sync.RWMutex
	port       int
	replId     string
	replOffset int
	replicaOf  ServerConfigReplicaOf
	RDB        ServerConfigRDB
}

func (server *Server) Configure() {
	replId, err := secureRandomString(40)
	if err == nil {
		server.config.replId = replId
	}

	server.config.mu.Lock()
	for i, arg := range os.Args {
		arg = strings.ReplaceAll(arg, "-", "")
		var val string

		if i+1 < len(os.Args) {
			val = os.Args[i+1]
		}

		switch arg {
		case "port":
			port, err := strconv.Atoi(val)
			if err == nil {
				server.config.port = port
			}

		case "dir":
			server.config.RDB.Dir = val

		case "replicaof":
			parts := strings.Split(val, " ")
			host := parts[0]
			portString := parts[1]
			port, err := strconv.Atoi(portString)

			if err == nil {
				server.config.replicaOf.host = host
				server.config.replicaOf.port = port

				conn, err := net.Dial("tcp", host+":"+portString)
				if err != nil {
					fmt.Println("Error connecting to master")
					break
				}
				fmt.Println("Connected to master")

				server.config.replicaOf.connection = Connection{
					id:         uuid.New(),
					role:       Slave,
					connection: conn,
					onCloseCh:  make(chan int),
				}
			}
		case "dbfilename":
			server.config.RDB.DBFileName = val
		}
	}
	server.config.mu.Unlock()
}
