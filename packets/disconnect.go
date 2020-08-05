package packets

import (
	"fmt"
	"io"
)

type DisconnectPacket struct {
	FixedHeader
}

var disconnectType = PacketType{
	name:     "DISCONNECT",
	packetId: 224,
}

// Name returns the packet type name.
func (d *DisconnectPacket) Name() string {
	return d.name
}

// CreatePacket creates a new packet with the appropriate FixedHeader.
// It sets default values where needed as well.
func (d *DisconnectPacket) CreatePacket() {
	d.FixedHeader = FixedHeader{PacketType: disconnectType, RemainingLength: 0}
}

func (d *DisconnectPacket) String() string {
	return fmt.Sprintf("%v", d.FixedHeader)
}

// Write creates the bytes.Buffer of the packet and writes them to
// the supplied io.Writer.
func (d *DisconnectPacket) Write(w io.Writer) error {
	packet := d.WriteHeader()
	_, err := packet.WriteTo(w)
	return err
}

// Read creates the packet from an io.Reader. It assumes that the
// first byte, the packet id, has already been read.
func (d *DisconnectPacket) Read(r io.Reader) error {
	var fh FixedHeader
	fh.PacketType = disconnectType
	err := fh.read(r)
	if err != nil {
		return err
	}
	d.FixedHeader = fh

	return nil
}
