package frame

import "fmt"

type Array struct {
	Length int
	Value  []Frame
}

func CreateArray(frames []Frame) *Array {
	return &Array{Length: len(frames), Value: frames}
}

func (array *Array) Encode() []byte {
	encoded := []byte(fmt.Sprintf("*%d\r\n", array.Length))
	for _, value := range array.Value {
		encoded = append(encoded, value.Encode()...)

	}

	return encoded
}

func (array *Array) Unpack() (command string, args []string, buf []byte) {
	command, _, _ = array.Value[0].Unpack()

	args = make([]string, array.Length-1)
	for i, value := range array.Value[1:] {
		arg, _, _ := value.Unpack()
		args[i] = arg
	}
	return command, args, nil
}
