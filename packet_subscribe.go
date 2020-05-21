package goqtt

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

func CreateSubscribePacket(topic string) (sp SubscribePacket) {
	sp.FixedHeader = FixedHeader{MessageType: "SUBSCRIBE"}
	sp.MessageId = []byte{0, 1}
	sp.Topics = []string{topic}
	sp.Qos = []byte{0}
	return
}

func (s *SubscribePacket) Write(w io.Writer, v bool) error {
	var body bytes.Buffer
	var err error

	body.Write(s.MessageId)
	for i, topic := range s.Topics {
		body.Write(encodeString(topic))
		body.WriteByte(s.Qos[i])
	}

	s.FixedHeader.RemainingLength = body.Len()
	packet := s.FixedHeader.WriteHeader()
	packet.Write(body.Bytes())

	if v {
		fmt.Println("BODY", body)
		fmt.Println("PACKET", packet)
	}
	_, err = packet.WriteTo(w)

	return err
}
