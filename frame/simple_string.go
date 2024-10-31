package frame

import (
	"fmt"
	"strings"
)

type SimpleString struct {
	Value string
}

func CreateSimpleString(value string) *SimpleString {
	return &SimpleString{Value: value}
}

func (simpleString *SimpleString) Encode() []byte {
	return []byte(fmt.Sprintf("+%s\r\n", simpleString.Value))
}

func (simpleString *SimpleString) Unpack() (command string, args []string, buf []byte) {
	parts := strings.Split(simpleString.Value, " ")
	return parts[0], parts[1:], nil
}
