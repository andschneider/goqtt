package packets

import (
	"fmt"
	"io"
)

type PingRespPacket struct {
	FixedHeader
}

func CreatePingRespPacket() (pp PingRespPacket) {
	pp.FixedHeader = FixedHeader{MessageType: "PINGRESP"}
	return
}

func (p *PingRespPacket) String() string {
	return fmt.Sprintf("%v", p.FixedHeader)
}

func (p *PingRespPacket) Write(w io.Writer) error {
	packet := p.WriteHeader()
	_, err := packet.WriteTo(w)
	return err
}

func (p *PingRespPacket) ReadPingRespPacket(r io.Reader) error {
	var fh FixedHeader
	fh.MessageType = "PINGRESP"
	err := fh.read(r)
	if err != nil {
		return err
	}
	p.FixedHeader = fh

	return nil
}
