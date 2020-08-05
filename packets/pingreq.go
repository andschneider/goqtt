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

// Name returns the packet type name.
func (pp *PingReqPacket) Name() string {
	return pp.name
}

// CreatePacket creates a new packet with the appropriate FixedHeader.
// It sets default values where needed as well.
func (pp *PingReqPacket) CreatePacket() {
	pp.FixedHeader = FixedHeader{PacketType: pingReqType, RemainingLength: 0}
}

func (pp *PingReqPacket) String() string {
	return fmt.Sprintf("%v", pp.FixedHeader)
}

// Write creates the bytes.Buffer of the packet and writes them to
// the supplied io.Writer.
func (pp *PingReqPacket) Write(w io.Writer) error {
	packet := pp.WriteHeader()
	_, err := packet.WriteTo(w)
	return err
}

// Read creates the packet from an io.Reader. It assumes that the
// first byte, the packet id, has already been read.
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
