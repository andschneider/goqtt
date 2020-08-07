package packets

import (
	"fmt"
	"io"
)

type PingRespPacket struct {
	FixedHeader
}

var pingRespType = PacketType{
	name:     "PINGRESP",
	packetId: 208,
}

// Name returns the packet type name.
func (p *PingRespPacket) Name() string {
	return p.name
}

// CreatePacket creates a new packet with the appropriate FixedHeader.
// It sets default values where needed as well.
func (p *PingRespPacket) CreatePacket() {
	p.FixedHeader = FixedHeader{PacketType: pingRespType}
}

func (p *PingRespPacket) String() string {
	return fmt.Sprintf("%v", p.FixedHeader)
}

// Write creates the bytes.Buffer of the packet and writes them to
// the supplied io.Writer.
func (p *PingRespPacket) Write(w io.Writer) error {
	packet := p.WriteHeader()
	_, err := packet.WriteTo(w)
	return err
}

// Read creates the packet from an io.Reader. It assumes that the
// first byte, the packet id, has already been read.
func (p *PingRespPacket) Read(r io.Reader) error {
	var fh FixedHeader
	fh.PacketType = pingRespType
	err := fh.read(r)
	if err != nil {
		return err
	}
	p.FixedHeader = fh

	return nil
}
