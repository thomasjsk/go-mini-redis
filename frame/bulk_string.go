package frame

import "fmt"

type BulkString struct {
	Length int
	Value  string
}

func CreateBulkString(value string) *BulkString {
	return &BulkString{Value: value, Length: len(value)}
}

func (bulkString *BulkString) Encode() []byte {
	stringVal := fmt.Sprintf("$-1\r\n")
	if len(bulkString.Value) > 0 {
		stringVal = fmt.Sprintf("$%d\r\n%s\r\n", bulkString.Length, bulkString.Value)
	}

	return []byte(stringVal)
}

func (bulkString *BulkString) Unpack() (command string, args []string, buf []byte) {
	return bulkString.Value, []string{bulkString.Value}, nil
}
