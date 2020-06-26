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

func TestReaderWithConnack(t *testing.T) {
	connack := bytes.NewBuffer(testConnackPacket)

	// read packet
	var ca ConnackPacket
	err := ca.Read(connack)
	if err != nil {
		t.Errorf("could not read %s packet: %v\n", ca.name, err)
	}

	var b bytes.Buffer
	err = ca.Write(&b)
	if err != nil {
		t.Errorf("could not write %s packet: %v\n", ca.name, err)
	}

	p, err := Reader(&b)
	if err != nil {
		t.Errorf("could not use Reader: %v\n", err)
	}
	fmt.Printf("packet type is: %T\n", p)

	if p != ca {
		t.Errorf("packets don't match:\n%+v\n%+v", p, ca)
	}

}
