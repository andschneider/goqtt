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
	name:     "UNSUBACK",
	packetId: 176,
}

// Name returns the packet type name.
func (ua *UnsubackPacket) Name() string {
	return ua.name
}

// CreatePacket creates a new packet with the appropriate FixedHeader.
// It sets default values where needed as well.
func (ua *UnsubackPacket) CreatePacket() {
	ua.FixedHeader = FixedHeader{PacketType: unsubackType}
	ua.MessageId = defaultMessageId
}

func (ua *UnsubackPacket) String() string {
	return fmt.Sprintf("%v messageid: %b", ua.FixedHeader, ua.MessageId)
}

// Write creates the bytes.Buffer of the packet and writes them to
// the supplied io.Writer.
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

// Read creates the packet from an io.Reader. It assumes that the
// first byte, the packet id, has already been read.
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
