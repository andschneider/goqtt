package packets

import (
	"bytes"
	"fmt"
	"testing"
)

const (
	testTopic = "test/topic"
)

func TestFixedHeader(t *testing.T) {
	var expected bytes.Buffer
	expected.Write([]byte{16, 0})

	fh := FixedHeader{PacketType: connectType}
	h := fh.WriteHeader()
	if !bytes.Equal(expected.Bytes(), h.Bytes()) {
		t.Errorf("Expected %b, got %b\n", expected, h)
	}
}

// Adds the packetId to the packet
var testFullConnackPacket = []byte{byte(32), byte(remainingLength), byte(sessionPresent), byte(returnCode)}

func TestReaderWithConnack(t *testing.T) {
	connack := bytes.NewBuffer(testConnackPacket)

	// read packet using the Read method directly
	var ca ConnackPacket
	err := ca.Read(connack)
	if err != nil {
		t.Errorf("could not read %s packet: %v\n", ca.name, err)
	}

	// Use the ReadPacket functon to return a Packet interface
	c2 := bytes.NewBuffer(testFullConnackPacket)
	p, err := ReadPacket(c2)
	if err != nil {
		t.Errorf("could not use ReadPacke: %v\n", err)
	}
	caNew := p.(*ConnackPacket)
	fmt.Printf("packet type is: %T\n", caNew)

	if *caNew != ca {
		t.Errorf("packets don't match:\n%+v\n%+v", *caNew, ca)
	}
}
