package packets

import (
	"fmt"
	"io"
)

type PingReqPacket struct {
	FixedHeader
}

var pingReqType = PacketType{
	name:     "PINGREQ",
	packetId: 192,
}

func (pp *PingReqPacket) Name() string {
	return pp.name
}

func (pp *PingReqPacket) CreatePacket() {
	pp.FixedHeader = FixedHeader{PacketType: pingReqType, RemainingLength: 0}
}

func (pp *PingReqPacket) String() string {
	return fmt.Sprintf("%v", pp.FixedHeader)
}

func (pp *PingReqPacket) Write(w io.Writer) error {
	packet := pp.WriteHeader()
	_, err := packet.WriteTo(w)
	return err
}

func (pp *PingReqPacket) Read(r io.Reader) error {
	var fh FixedHeader
	fh.PacketType = pingReqType
	err := fh.read(r)
	if err != nil {
		return err
	}
	pp.FixedHeader = fh

	return nil
}
