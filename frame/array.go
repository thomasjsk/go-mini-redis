package frame

type Array struct {
	Length int
	Value  []Frame
}

func (array *Array) Encode() []byte {
	return []byte("")
}

func (array *Array) Unpack() (command string, args []string) {
	command, _ = array.Value[0].Unpack()

	args = make([]string, array.Length-1)
	for i, value := range array.Value[1:] {
		arg, _ := value.Unpack()
		args[i] = arg
	}
	return command, args
}
