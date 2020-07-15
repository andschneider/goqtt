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

func (d *DisconnectPacket) Name() string {
	return d.name
}

func (d *DisconnectPacket) CreatePacket() {
	d.FixedHeader = FixedHeader{PacketType: disconnectType, RemainingLength: 0}
}

func (d *DisconnectPacket) String() string {
	return fmt.Sprintf("%v", d.FixedHeader)
}

func (d *DisconnectPacket) Write(w io.Writer) error {
	packet := d.WriteHeader()
	_, err := packet.WriteTo(w)
	return err
}

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
