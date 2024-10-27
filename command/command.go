package command

import (
	"github.com/thomasjsk/go-mini-redis/frame"
)

type Command struct {
	Name string
	Args []string
}

func (cmd Command) FromFrame(f *frame.Frame) Command {
	switch f := (*f).(type) {
	case *frame.Array:
		name, args := f.Unpack()
		cmd.Name = name
		cmd.Args = args
		break
	default:
		break
	}

	return cmd
}
