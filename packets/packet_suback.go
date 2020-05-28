package packets

import (
	"fmt"
	"io"
)

type SubackPacket struct {
	FixedHeader
	MessageId   uint16
	ReturnCodes []byte
}

var subackType = PacketType{
	name:     "SUBACK",
	packetId: 144,
}

func (sa *SubackPacket) String() string {
	return fmt.Sprintf("%v messageid: %d", sa.FixedHeader, sa.MessageId)
}

func (sa *SubackPacket) ReadSubackPacket(r io.Reader) error {
	var fh FixedHeader
	fh.PacketType = subackType
	err := fh.read(r)
	if err != nil {
		return err
	}
	sa.FixedHeader = fh

	sa.MessageId, err = decodeUint16(r)
	if err != nil {
		return err
	}

	// TODO this only works if the suback has a single topic. need to expand to a list of topics
	rc, err := decodeByte(r)
	if err != nil {
		return err
	}
	sa.ReturnCodes = []byte{rc}

	return nil
}
