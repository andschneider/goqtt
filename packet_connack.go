package goqtt

import (
	"fmt"
	"io"
)

type ConnackPacket struct {
	FixedHeader
	SessionPresent byte
	ReturnCode     byte
}

func (c *ConnackPacket) String() string {
	return fmt.Sprintf("%v sessionpresent: %d returncode: %d", c.FixedHeader, c.SessionPresent, c.ReturnCode)
}

func (c *ConnackPacket) ReadConnackPacket(r io.Reader) error {
	var fh FixedHeader
	fh.MessageType = "CONNACK"
	err := fh.read(r)
	if err != nil {
		return err
	}
	c.FixedHeader = fh

	sp, err := decodeByte(r)
	if err != nil {
		return err
	}
	c.SessionPresent = sp

	rc, err := decodeByte(r)
	if err != nil {
		return err
	}
	c.ReturnCode = rc

	return nil
}
