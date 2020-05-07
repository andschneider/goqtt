package main

import (
	"bytes"
	"encoding/binary"
)

func (fh *FixedHeader) WriteHeader() bytes.Buffer {
	var header bytes.Buffer
	header.WriteByte(fh.MessageType)
	header.Write(encodeLength(fh.RemainingLength))
	return header
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

func encodeClientId(ci string) []byte {
	cb := []byte(ci)
	l := make([]byte, 2)
	binary.BigEndian.PutUint16(l, uint16(len(cb)))
	return append(l, cb...)
}
