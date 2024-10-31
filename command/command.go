package command

import (
	"fmt"
	"github.com/thomasjsk/go-mini-redis/frame"
)

type Command struct {
	Name string
	Args []string
	Buf  []byte
}

func (cmd Command) FromFrame(f *frame.Frame) Command {
	switch f := (*f).(type) {
	case *frame.Array:
		name, args, _ := f.Unpack()
		cmd.Name = name
		cmd.Args = args
		break
	case *frame.BytesFrame:
		name, _, buf := f.Unpack()
		cmd.Name = name
		cmd.Buf = buf
	case *frame.SimpleString:
		name, _, _ := f.Unpack()
		cmd.Name = name
	default:
		fmt.Println("Unsupported command type", f)
		break
	}

	return cmd
}
