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

var subscribeType = PacketType{
	name:     "SUBSCRIBE",
	packetId: 130,
}

// Name returns the packet type name.
func (s *SubscribePacket) Name() string {
	return s.name
}

// CreateSubscribePacket wraps CreatePacket and adds a topic.
func (s *SubscribePacket) CreateSubscribePacket(topic string) {
	s.CreatePacket()
	s.Topics = []string{topic}
}

// CreatePacket creates a new packet with the appropriate FixedHeader.
// It sets default values where needed as well.
func (s *SubscribePacket) CreatePacket() {
	s.FixedHeader = FixedHeader{PacketType: subscribeType}
	s.MessageId = defaultMessageId
	s.Qos = []byte{0}
}

func (s *SubscribePacket) String() string {
	return fmt.Sprintf("%v messageid: %v topics: %s", s.FixedHeader, s.MessageId, s.Topics)
}

// Write creates the bytes.Buffer of the packet and writes them to
// the supplied io.Writer.
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

// Read creates the packet from an io.Reader. It assumes that the
// first byte, the packet id, has already been read.
func (s *SubscribePacket) Read(r io.Reader) error {
	var fh FixedHeader
	fh.PacketType = subscribeType
	err := fh.read(r)
	if err != nil {
		return fmt.Errorf("could not read in header: %v", err)
	}

	s.FixedHeader = fh
	s.MessageId, err = decodeMessageId(r)
	if err != nil {
		return err
	}
	payloadLength := s.RemainingLength - 2
	for payloadLength > 0 {
		topic, err := decodeString(r)
		if err != nil {
			return err
		}
		s.Topics = append(s.Topics, topic)
		qos, err := decodeByte(r)
		if err != nil {
			return err
		}
		s.Qos = append(s.Qos, qos)
		payloadLength -= 2 + len(topic) + 1 //2 bytes of string length, plus string, plus 1 byte for Qos
	}
	return nil
}
