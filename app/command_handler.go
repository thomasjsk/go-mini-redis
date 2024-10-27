package main

import (
	"fmt"
	"github.com/thomasjsk/go-mini-redis/command"
	"github.com/thomasjsk/go-mini-redis/frame"
	"strconv"
	"time"
)

func HandleCommand(command command.Command) frame.Frame {
	switch command.Name {
	case "ping":
		return ping()
	case "echo":
		return echo(command.Args)
	case "set":
		return set(command.Args)
	case "get":
		return get(command.Args[0])
	default:
		return ping()
	}

}

func ping() frame.Frame {
	return &frame.SimpleString{
		Value: "PONG",
	}
}

func echo(args []string) frame.Frame {
	return &frame.SimpleString{Value: args[0]}
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

	SetMutex.Lock()
	SetStorage[k] = value
	SetMutex.Unlock()

	return &frame.SimpleString{Value: "OK"}
}

func get(key string) frame.Frame {
	SetMutex.Lock()
	storageValue, ok := SetStorage[key]
	SetMutex.Unlock()

	if !ok {
		return &frame.BulkString{
			Value: "",
		}
	}

	if !storageValue.ExpireAt.IsZero() && time.Now().After(storageValue.ExpireAt) {
		SetMutex.Lock()
		delete(SetStorage, key)
		SetMutex.Unlock()

		return &frame.BulkString{
			Value: "",
		}
	}

	fmt.Println(storageValue.Value, storageValue.ExpireAt)

	return &frame.SimpleString{Value: storageValue.Value}
}
