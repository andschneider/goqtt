package packets

import (
	"bytes"
	"fmt"
	"testing"
)

var testUnsubackPacket = []byte{2, 0, 1, 0}

func TestUnsubackPacket(t *testing.T) {
	unsuback := bytes.NewBuffer(testUnsubackPacket)

	ua := UnsubackPacket{}
	err := ua.Read(unsuback)
	if err != nil {
		t.Errorf("could not read suback packet: %v\n", err)
	}
	fmt.Printf("unsuback packet: %s\n", ua.String())
}

func TestUnsubackPacket_Write(t *testing.T) {
	var buf bytes.Buffer
	var uaRead, uaWrite UnsubackPacket
	suback := bytes.NewBuffer(testUnsubackPacket)

	err := uaRead.Read(suback)
	if err != nil {
		t.Errorf("could not read %s packet: %v\n", uaRead.Name, err)
	}
	fmt.Printf("read %s packet: %+v\n", uaRead.Name, uaRead)

	// create packet
	uaWrite.CreatePacket()
	err = uaWrite.Write(&buf)
	if err != nil {
		t.Errorf("could not %s suback packet: %v", uaWrite.Name, err)
	}
	fmt.Printf("write %s packet: %+v\n", uaWrite.Name, uaWrite)

	// verify they match
	if !bytes.Equal(uaRead.MessageId, uaWrite.MessageId) {
		t.Errorf("unsuback messageids don't match")
	}
	if uaRead.FixedHeader != uaWrite.FixedHeader {
		t.Errorf("unsuback fixedheaders don't match")
	}
}
