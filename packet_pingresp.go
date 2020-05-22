package goqtt

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

func (p *PingRespPacket) Write(w io.Writer, v bool) error {
	packet := p.FixedHeader.WriteHeader()
	if v {
		fmt.Println("PACKET", packet)
	}
	_, err := packet.WriteTo(w)
	return err
}
