package packets

import (
	"fmt"
	"io"
)

type PingReqPacket struct {
	FixedHeader
}

func (pp *PingReqPacket) String() string {
	return fmt.Sprintf("%v", pp.FixedHeader)
}

func CreatePingReqPacket() (pp PingReqPacket) {
	pp.FixedHeader = FixedHeader{MessageType: "PINGREQ", RemainingLength: 0}
	return
}

func (pp *PingReqPacket) Write(w io.Writer, v bool) error {
	packet := pp.FixedHeader.WriteHeader()
	if v {
		fmt.Println("PACKET", packet)
	}
	_, err := packet.WriteTo(w)
	return err
}
