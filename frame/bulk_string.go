package frame

import "fmt"

type BulkString struct {
	Length int
	Value  string
}

func (bulkString *BulkString) Encode() []byte {
	stringVal := fmt.Sprintf("$-1\r\n")
	if len(bulkString.Value) > 0 {
		stringVal = fmt.Sprintf("$%d\r\n%s\r\n", bulkString.Length, bulkString.Value)
	}

	return []byte(stringVal)
}

func (bulkString *BulkString) Unpack() (command string, args []string) {
	return bulkString.Value, nil
}
