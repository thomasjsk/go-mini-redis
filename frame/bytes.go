package frame

type BytesFrame struct {
	Value []byte
}

func (frame *BytesFrame) Encode() []byte {
	return frame.Value
}

func (frame *BytesFrame) Unpack() (command string, args []string, buf []byte) {
	return "BTYES", []string{string(frame.Value)}, frame.Value
}
