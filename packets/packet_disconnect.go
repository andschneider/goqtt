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

func (d *DisconnectPacket) String() string {
	return fmt.Sprintf("%v", d.FixedHeader)
}

func CreateDisconnectPacket() (d DisconnectPacket) {
	d.FixedHeader = FixedHeader{PacketType: disconnectType, RemainingLength: 0}
	return
}

func (d *DisconnectPacket) Write(w io.Writer) error {
	packet := d.WriteHeader()
	_, err := packet.WriteTo(w)
	return err
}

func (d *DisconnectPacket) ReadDisconnectPacket(r io.Reader) error {
	var fh FixedHeader
	fh.PacketType = disconnectType
	err := fh.read(r)
	if err != nil {
		return err
	}
	d.FixedHeader = fh

	return nil
}
