package packets

import (
	"bytes"
	"fmt"
	"testing"
)

var testSubackPacket = []byte{3, 0, 1, 0}

func TestSubackPacket(t *testing.T) {
	suback := bytes.NewBuffer(testSubackPacket)

	sa := SubackPacket{}
	err := sa.Read(suback)
	if err != nil {
		t.Errorf("could not read suback packet: %v\n", err)
	}
	fmt.Printf("%s packet: %+v\n", sa.name, sa)
}

func TestSubackPacket_Write(t *testing.T) {
	var buf bytes.Buffer
	var saRead, saWrite SubackPacket
	suback := bytes.NewBuffer(testSubackPacket)

	err := saRead.Read(suback)
	if err != nil {
		t.Errorf("could not read %s packet: %v\n", saRead.name, err)
	}
	fmt.Printf("read %s packet: %+v\n", saRead.name, saRead)

	// create packet
	saWrite.CreatePacket()
	err = saWrite.Write(&buf)
	if err != nil {
		t.Errorf("could not %s suback packet: %v", saWrite.name, err)
	}
	fmt.Printf("write %s packet: %+v\n", saWrite.name, saWrite)

	// verify they match
	if !bytes.Equal(saRead.MessageId, saWrite.MessageId) {
		t.Errorf("suback messageids don't match")
	}
	if !bytes.Equal(saRead.ReturnCodes, saWrite.ReturnCodes) {
		t.Errorf("suback returncodes don't match")
	}
	if saRead.FixedHeader != saWrite.FixedHeader {
		t.Errorf("suback fixedheaders don't match")
	}
}