package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	MQTT5 = 5 // 0000101 in binary
	MQTT3 = 4 // 0000100 in binary
)

var MessageTypes = map[string]byte{
	"CONNECT":   16,  // 00010000 in binary
	"SUBSCRIBE": 130, // 10000010 in binary
}

type FixedHeader struct {
	MessageType     string
	RemainingLength int
}

func (fh *FixedHeader) WriteHeader() (header bytes.Buffer) {
	t, ok := MessageTypes[fh.MessageType]
	if !ok {
		fmt.Println("wrong message type, must be ...") // TODO: make this better
	}
	header.WriteByte(t)
	header.Write(encodeLength(fh.RemainingLength))
	return
}

func encodeLength(length int) []byte {
	var encLength []byte
	for {
		digit := byte(length % 128)
		length /= 128
		if length > 0 {
			digit |= 0x80
		}
		encLength = append(encLength, digit)
		if length == 0 {
			break
		}
	}
	return encLength
}

func encodeString(s string) []byte {
	b := []byte(s)
	e := make([]byte, 2)
	binary.BigEndian.PutUint16(e, uint16(len(b)))
	return append(e, b...)
}
