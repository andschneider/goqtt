package packets

import (
	"fmt"
	"io"
)

type PingRespPacket struct {
	FixedHeader
}

var pingRespType = PacketType{
	Name:     "PINGRESP",
	packetId: 208,
}

func (p *PingRespPacket) CreatePacket() {
	p.FixedHeader = FixedHeader{PacketType: pingRespType}
}

func (p *PingRespPacket) String() string {
	return fmt.Sprintf("%v", p.FixedHeader)
}

func (p *PingRespPacket) Write(w io.Writer) error {
	packet := p.WriteHeader()
	_, err := packet.WriteTo(w)
	return err
}

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
