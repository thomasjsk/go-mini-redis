package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/thomasjsk/go-mini-redis/frame"
	"net"
)

type Role string

const (
	Master Role = "MASTER"
	Slave       = "SLAVE"
)

type Connection struct {
	id         uuid.UUID
	role       Role
	connection net.Conn
	onCloseCh  chan int
}

func (c *Connection) Close() {
	err := c.connection.Close()
	if err != nil {
		fmt.Println("Error closing connection")
	}
}

func (c *Connection) Read(buf []byte) (int, error) {
	return c.connection.Read(buf)
}

func (c *Connection) Write(f frame.Frame) (int, error) {
	return c.connection.Write(f.Encode())
}

func (c *Connection) WriteRaw(bytes []byte) (int, error) {
	return c.connection.Write(bytes)
}
