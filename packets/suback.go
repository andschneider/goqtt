package packets

import (
	"bytes"
	"fmt"
	"io"
)

type SubackPacket struct {
	FixedHeader
	MessageId   []byte
	ReturnCodes []byte
}

var subackType = PacketType{
	name:     "SUBACK",
	packetId: 144,
}

// Name returns the packet type name.
func (sa *SubackPacket) Name() string {
	return sa.name
}

// CreatePacket creates a new packet with the appropriate FixedHeader.
// It sets default values where needed as well.
func (sa *SubackPacket) CreatePacket() {
	sa.FixedHeader = FixedHeader{PacketType: subackType}
	sa.MessageId = []byte{0, 1}
	// TODO expand to more than one topic
	sa.ReturnCodes = []byte{0}
}

func (sa *SubackPacket) String() string {
	return fmt.Sprintf("%v messageid: %b returncodes: %b", sa.FixedHeader, sa.MessageId, sa.ReturnCodes)
}

// Write creates the bytes.Buffer of the packet and writes them to
// the supplied io.Writer.
func (sa *SubackPacket) Write(w io.Writer) error {
	var body bytes.Buffer
	var err error

	body.Write(sa.MessageId)
	body.Write(sa.ReturnCodes)

	sa.RemainingLength = body.Len()
	packet := sa.WriteHeader()
	packet.Write(body.Bytes())
	_, err = packet.WriteTo(w)

	return err
}

// Read creates the packet from an io.Reader. It assumes that the
// first byte, the packet id, has already been read.
func (sa *SubackPacket) Read(r io.Reader) error {
	var fh FixedHeader
	fh.PacketType = subackType
	err := fh.read(r)
	if err != nil {
		return err
	}
	sa.FixedHeader = fh

	sa.MessageId, err = decodeMessageId(r)
	if err != nil {
		return err
	}

	// TODO this only works if the suback has a single topic. need to expand to a list of topics
	rc, err := decodeByte(r)
	if err != nil {
		return err
	}
	sa.ReturnCodes = []byte{rc}

	return nil
}
