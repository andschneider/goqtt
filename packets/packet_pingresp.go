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

func (p *PingRespPacket) String() string {
	return fmt.Sprintf("%v", p.FixedHeader)
}

func CreatePingRespPacket() (pp PingRespPacket) {
	pp.FixedHeader = FixedHeader{PacketType: pingRespType}
	return
}

func (p *PingRespPacket) Write(w io.Writer) error {
	packet := p.WriteHeader()
	_, err := packet.WriteTo(w)
	return err
}

func (p *PingRespPacket) ReadPingRespPacket(r io.Reader) error {
	var fh FixedHeader
	fh.PacketType = pingRespType
	err := fh.read(r)
	if err != nil {
		return err
	}
	p.FixedHeader = fh

	return nil
}
