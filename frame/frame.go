package frame

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Frame interface {
	Encode() []byte
	Unpack() (command string, args []string)
}

type RawFrame []byte

func (rawFrame *RawFrame) Decode(cursor *int) (frame Frame, err error) {
	frameType := (*rawFrame).readType(cursor)

	switch frameType {
	case "+":
		return rawFrame.DecodeSimpleString(cursor)
	case "$":
		return rawFrame.DecodeBulkString(cursor)
	case "*":
		return rawFrame.DecodeArray(cursor)
	default:
		return nil, errors.New("Unsupported frame type: " + frameType)
	}
}

func (rawFrame *RawFrame) DecodeSimpleString(cursor *int) (frame *SimpleString, err error) {
	line, err := rawFrame.readLine(cursor)

	if err != nil {
		return &SimpleString{}, err
	}

	return &SimpleString{Value: strings.ToLower(line)}, nil
}

func (rawFrame *RawFrame) DecodeBulkString(cursor *int) (frame *BulkString, err error) {
	line, err := rawFrame.readLine(cursor)

	if err != nil {
		return &BulkString{}, err
	}

	strLen, err := strconv.Atoi(string(line[1]))
	if err != nil {
		return &BulkString{}, errors.New("invalid bulk string length: " + line)
	}

	line, err = rawFrame.readLine(cursor)

	return &BulkString{Value: strings.ToLower(line), Length: strLen}, nil
}

func (rawFrame *RawFrame) DecodeArray(cursor *int) (frame *Array, err error) {
	line, err := rawFrame.readLine(cursor)

	if err != nil {
		return &Array{}, err
	}

	arrLen, err := strconv.Atoi(string(line[1]))
	array := Array{Length: arrLen, Value: make([]Frame, arrLen)}
	for i := range array.Length {
		frame, err := rawFrame.Decode(cursor)
		if err != nil {
			fmt.Println("error decoding array item")
		}

		array.Value[i] = frame

	}

	return &array, nil
}
func (rawFrame *RawFrame) readType(cursor *int) string {
	frameType := string((*rawFrame)[*cursor])

	return frameType
}

func (rawFrame *RawFrame) readLine(cursor *int) (line string, err error) {
	for _i := range (*rawFrame)[*cursor:] {
		i := *cursor + _i

		if _i > 1 && rune((*rawFrame)[i-1]) == '\r' && rune((*rawFrame)[i]) == '\n' {
			line := string((*rawFrame)[*cursor : i-1])
			*cursor = i + 1

			return line, nil
		}
	}

	return "", errors.New("cannot read frame")
}
