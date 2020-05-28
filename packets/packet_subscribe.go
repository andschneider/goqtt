package packets

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

func (s *SubscribePacket) String() string {
	return fmt.Sprintf("%v messageid: %v topics: %s", s.FixedHeader, s.MessageId, s.Topics)
}

func CreateSubscribePacket(topic string) (sp SubscribePacket) {
	sp.FixedHeader = FixedHeader{MessageType: "SUBSCRIBE"}
	sp.MessageId = []byte{0, 1}
	sp.Topics = []string{topic}
	sp.Qos = []byte{0}
	return
}

func (s *SubscribePacket) Write(w io.Writer) error {
	var body bytes.Buffer
	var err error

	body.Write(s.MessageId)
	for i, topic := range s.Topics {
		body.Write(encodeString(topic))
		body.WriteByte(s.Qos[i])
	}

	s.RemainingLength = body.Len()
	packet := s.WriteHeader()
	packet.Write(body.Bytes())
	_, err = packet.WriteTo(w)

	return err
}
