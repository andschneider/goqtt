package packets

import (
	"bytes"
	"fmt"
	"io"
)

type UnsubscribePacket struct {
	FixedHeader
	MessageId []byte
	Topics    []string
}

var unsubscribeType = PacketType{
	name:     "UNSUBSCRIBE",
	packetId: 162,
}

func (up *UnsubscribePacket) String() string {
	return fmt.Sprintf("%v messageid: %b topics %s", up.FixedHeader, up.MessageId, up.Topics)
}

func CreateUnsubscribePacket(topic string) (up UnsubscribePacket) {
	up.FixedHeader = FixedHeader{PacketType: unsubscribeType}
	up.MessageId = []byte{0, 1}
	up.Topics = []string{topic}
	return
}

func (up *UnsubscribePacket) Write(w io.Writer) error {
	var body bytes.Buffer
	var err error

	body.Write(up.MessageId)
	for _, topic := range up.Topics {
		body.Write(encodeString(topic))
	}

	up.RemainingLength = body.Len()
	packet := up.WriteHeader()
	packet.Write(body.Bytes())
	_, err = packet.WriteTo(w)

	return err
}

func (up *UnsubscribePacket) Read(r io.Reader) error {
	var fh FixedHeader
	fh.PacketType = unsubscribeType
	err := fh.read(r)
	if err != nil {
		return err
	}
	up.FixedHeader = fh

	up.MessageId, err = decodeMessageId(r)
	if err != nil {
		return err
	}
	payloadLength := up.RemainingLength - 2
	for payloadLength > 0 {
		topic, err := decodeString(r)
		if err != nil {
			return err
		}
		up.Topics = append(up.Topics, topic)
		payloadLength -= 2 + len(topic) //2 bytes of string length plus string
	}
	return nil
}
