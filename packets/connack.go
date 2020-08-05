package packets

import (
	"bytes"
	"fmt"
	"io"
)

type ConnackPacket struct {
	FixedHeader
	SessionPresent byte
	ReturnCode     byte
}

var connackType = PacketType{
	name:     "CONNACK",
	packetId: 32,
}

// Name returns the packet type name.
func (c *ConnackPacket) Name() string {
	return c.name
}

// CreatePacket creates a new packet with the appropriate FixedHeader.
// It sets default values where needed as well.
func (c *ConnackPacket) CreatePacket() {
	c.FixedHeader = FixedHeader{PacketType: connackType, RemainingLength: 2}
	// hardcode for a connect packet with connect flags of 00000010
	c.SessionPresent = 0
	c.ReturnCode = 0
}

func (c *ConnackPacket) String() string {
	return fmt.Sprintf("%v sessionpresent: %d returncode: %d", c.FixedHeader, c.SessionPresent, c.ReturnCode)
}

// Write creates the bytes.Buffer of the packet and writes them to
// the supplied io.Writer.
func (c *ConnackPacket) Write(w io.Writer) error {
	var body bytes.Buffer
	body.WriteByte(c.SessionPresent)
	body.WriteByte(c.ReturnCode)

	if body.Len() != c.RemainingLength {
		return fmt.Errorf("body of CONNACK is incorrect length: %d instead of %d", body.Len(), c.RemainingLength)
	}
	packet := c.WriteHeader()
	packet.Write(body.Bytes())
	_, err := packet.WriteTo(w)
	return err
}

// Read creates the packet from an io.Reader. It assumes that the
// first byte, the packet id, has already been read.
func (c *ConnackPacket) Read(r io.Reader) error {
	var fh FixedHeader
	fh.PacketType = connackType
	err := fh.read(r)
	if err != nil {
		return err
	}
	c.FixedHeader = fh

	sp, err := decodeByte(r)
	if err != nil {
		return err
	}
	c.SessionPresent = sp

	rc, err := decodeByte(r)
	if err != nil {
		return err
	}
	c.ReturnCode = rc

	return nil
}
