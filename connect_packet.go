package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	MQTT5       = 5  // 0000101 in binary
	MQTT3       = 4  // 0000100 in binary
	MessageType = 16 // 0010000 in binary
)

type ConnectPacket struct {
	FixedHeader
	ProtocolName    []byte
	ProtocolVersion byte
	ConnectFlags    byte
	KeepAlive       []byte

	ClientIdentifier string
}

type FixedHeader struct {
	MessageType     byte
	RemainingLength int
}

func create() ConnectPacket {
	fh := FixedHeader{MessageType: MessageType}
	cp := ConnectPacket{FixedHeader: fh}

	cp.ProtocolName = []byte{0, 4, 77, 81, 84, 84} // "04MQTT"
	cp.ProtocolVersion = MQTT5
	cp.ConnectFlags = 2
	cp.KeepAlive = []byte{0, 60}
	cp.ClientIdentifier = "andr"

	return cp
}

func (fh *FixedHeader) WriteHeader() bytes.Buffer {
	var header bytes.Buffer
	header.WriteByte(fh.MessageType)
	header.Write(encodeLength(fh.RemainingLength))
	return header
}

func (c *ConnectPacket) Write(w io.Writer) error {
	var body bytes.Buffer
	var err error

	fmt.Println("BODY", body)
	body.Write(c.ProtocolName)
	body.WriteByte(c.ProtocolVersion)
	body.WriteByte(c.ConnectFlags)
	body.Write(c.KeepAlive)
	body.Write(encodeClientId(c.ClientIdentifier))

	fmt.Println("BODY", body)

	c.FixedHeader.RemainingLength = body.Len()

	packet := c.FixedHeader.WriteHeader()
	packet.Write(body.Bytes())
	fmt.Println("PACKET", packet)
	_, err = packet.WriteTo(w)

	return err
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
	test := append(l, cb...)
	fmt.Println("CID", test)
	return test
}
