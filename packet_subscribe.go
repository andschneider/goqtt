package main

import (
	"bytes"
	"fmt"
	"io"
)

type SubscribePacket struct {
	FixedHeader
	MessageId []byte
	Topics    []string
	Qos       []byte
}

func createSubPacket() SubscribePacket {
	fh := FixedHeader{MessageType: 128}
	sp := SubscribePacket{FixedHeader: fh}

	sp.MessageId = []byte{0, 1}
	sp.Topics = []string{"test/topic"}
	sp.Qos = []byte{0}

	return sp
}

func (s *SubscribePacket) Write(w io.Writer) error {
	var body bytes.Buffer
	var err error

	body.Write(s.MessageId)
	for i, topic := range s.Topics {
		body.Write(encodeClientId(topic))
		body.WriteByte(s.Qos[i])
	}
	fmt.Println("BODY", body)

	s.FixedHeader.RemainingLength = body.Len()
	packet := s.FixedHeader.WriteHeader()
	packet.Write(body.Bytes())
	fmt.Println("PACKET", packet)
	_, err = packet.WriteTo(w)

	return err
}
