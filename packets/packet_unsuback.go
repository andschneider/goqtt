package packets

import (
	"bytes"
	"fmt"
	"io"
)

type UnsubackPacket struct {
	FixedHeader
	MessageId []byte
}

var unsubackType = PacketType{
	Name:     "UNSUBACK",
	packetId: 176,
}

func (ua *UnsubackPacket) CreatePacket() {
	ua.FixedHeader = FixedHeader{PacketType: unsubackType}
	ua.MessageId = []byte{0, 1}
}

func (ua *UnsubackPacket) String() string {
	return fmt.Sprintf("%v messageid: %b", ua.FixedHeader, ua.MessageId)
}

func (ua *UnsubackPacket) Write(w io.Writer) error {
	var body bytes.Buffer
	var err error

	body.Write(ua.MessageId)

	ua.RemainingLength = body.Len()
	packet := ua.WriteHeader()
	packet.Write(body.Bytes())
	_, err = packet.WriteTo(w)

	return err
}

func (ua *UnsubackPacket) Read(r io.Reader) error {
	var fh FixedHeader
	fh.PacketType = unsubackType
	err := fh.read(r)
	if err != nil {
		return err
	}
	ua.FixedHeader = fh

	ua.MessageId, err = decodeMessageId(r)
	if err != nil {
		return err
	}

	return nil
}
