package main

import (
	"encoding/hex"
	"fmt"
	"github.com/thomasjsk/go-mini-redis/command"
	"github.com/thomasjsk/go-mini-redis/frame"
	"strconv"
	"time"
)

func HandleCommand(cmd command.Command, connection *Connection, server *Server, slave bool) (err error) {
	fmt.Println("HandleCommand.Name", cmd.Name, server.config.port)
	var response frame.Frame
	switch cmd.Name {
	case "ping":
		response = ping()
	case "echo":
		response = echo(cmd.Args)
	case "set":
		response = set(cmd.Args)
	case "get":
		response = get(cmd.Args[0])
	case "info":
		response = info(cmd.Args[0], &server.config)
	case "config":
		response = config(cmd.Args, &server.config)
	case "replconf":
		response = replconf(cmd.Args)
	case "psync":
		return psync(cmd.Args, connection, &server.config)
	case "fullresync":
		return fullresync(cmd)
	default:
		fmt.Println("Unsupported cmd:", cmd)
		response = nil
	}

	if response == nil {
		return nil
	}

	mode := "MASTER"
	if slave {
		mode = "SLAVE"
	}
	if !slave {
		fmt.Printf("[%s] responds with %s; \n", mode, response)
		_, err = connection.Write(response)
		if err != nil {
			fmt.Println("Error while writing the response", err)

			return err
		}
	} else {
		fmt.Printf("Wanted to respond with %s, but it's a slave...\n", response)
	}

	return nil
}

func ping() frame.Frame {
	return frame.CreateSimpleString("PONG")
}

func echo(args []string) frame.Frame {
	return frame.CreateSimpleString(args[0])
}

func set(args []string) frame.Frame {
	k := args[0]
	v := args[1]

	value := StorageValue{
		Value: v,
	}
	if len(args) > 2 {
		nextArg := args[2]

		if nextArg == "px" && len(args) > 3 {
			expiration := args[3]

			if len(expiration) > 0 {
				timestamp, err := strconv.ParseInt(expiration, 10, 64)
				if err == nil {
					value.ExpireAt = time.Now().Add(time.Duration(timestamp * 1000000))
				}

			}

		}
	}

	StorageMutex.Lock()
	Storage[k] = value
	StorageMutex.Unlock()

	return frame.CreateSimpleString("OK")
}

func get(key string) frame.Frame {
	StorageMutex.Lock()
	storageValue, ok := Storage[key]
	StorageMutex.Unlock()

	if !ok {
		return frame.CreateBulkString("")
	}

	if !storageValue.ExpireAt.IsZero() && time.Now().After(storageValue.ExpireAt) {
		StorageMutex.Lock()
		delete(Storage, key)
		StorageMutex.Unlock()

		return frame.CreateBulkString("")
	}

	fmt.Println(storageValue.Value, storageValue.ExpireAt)

	return frame.CreateSimpleString(storageValue.Value)
}

func info(arg string, serverConfig *ServerConfig) frame.Frame {
	buildRecord := func(k string, v string) string {
		return k + ":" + v + "\r\n"
	}

	if arg == "replication" {
		var result string

		if len(serverConfig.replicaOf.host) > 0 && serverConfig.replicaOf.port > 0 {
			result += buildRecord("role", "slave")
		} else {
			result += buildRecord("role", "master")
		}

		result += buildRecord("master_replid", serverConfig.replId)
		result += buildRecord("master_repl_offset", strconv.Itoa(serverConfig.replOffset))
		return frame.CreateBulkString(result)
	}

	return frame.CreateBulkString("")

}

func config(args []string, serverConfig *ServerConfig) frame.Frame {
	tail := args[0]
	arg := args[1]
	var value string

	if tail != "get" {
		return frame.CreateBulkString("")
	}

	switch arg {
	case "dir":
		value = serverConfig.RDB.Dir
	case "dbfilename":
		value = serverConfig.RDB.DBFileName
	}

	return frame.CreateArray([]frame.Frame{
		frame.CreateBulkString(arg),
		frame.CreateBulkString(value),
	})
}

func replconf(args []string) frame.Frame {
	param := args[0]
	value := args[1]

	if param == "listening-port" && len(value) > 0 {
		return frame.CreateSimpleString("OK")
	}

	if param == "capa" && value == "psync2" {
		return frame.CreateSimpleString("OK")
	}

	return frame.CreateSimpleString("")
}

func psync(args []string, conn *Connection, config *ServerConfig) error {

	reqReplId := args[0]

	var replId string
	var offset string
	if reqReplId == "?" {
		replId = config.replId
		config.mu.Lock()
		config.mu.Unlock()
		offset = strconv.Itoa(config.replOffset)
	}

	conn.role = Slave

	rdb, err := hex.DecodeString("524544495330303131fa0972656469732d76657205372e322e30fa0a72656469732d62697473c040fa056374696d65c26d08bc65fa08757365642d6d656dc2b0c41000fa08616f662d62617365c000fff06e3bfec0ff5aa2")
	if err != nil {
		return err
	}

	_, err = conn.Write(frame.CreateSimpleString(fmt.Sprintf("FULLRESYNC %s %s", replId, offset)))
	if err != nil {
		fmt.Println("Error while writing the response", err)
	}
	_, err = conn.WriteRaw(append([]byte(fmt.Sprintf("$%d\r\n", len(rdb))), rdb...))
	if err != nil {
		fmt.Println("Error while writing the response", err)
	}

	return nil

}

func fullresync(_ command.Command) error {
	return nil
}
