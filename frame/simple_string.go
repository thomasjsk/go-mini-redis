package frame

import "fmt"

type SimpleString struct {
	Value string
}

func (simpleString *SimpleString) Encode() []byte {
	return []byte(fmt.Sprintf("+%s\r\n", simpleString.Value))
}

func (simpleString *SimpleString) Unpack() (command string, args []string) {
	return simpleString.Value, nil
}
