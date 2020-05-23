package goqtt

import (
	"fmt"
	"io"
)

type PingReqPacket struct {
	FixedHeader
}

func CreatePingReqPacket() (pp PingReqPacket) {
	pp.FixedHeader = FixedHeader{MessageType: "PINGREQ", RemainingLength: 0}
	return
}

func (p *PingReqPacket) Write(w io.Writer, v bool) error {
	packet := p.FixedHeader.WriteHeader()
	if v {
		fmt.Println("PACKET", packet)
	}
	_, err := packet.WriteTo(w)
	return err
}
