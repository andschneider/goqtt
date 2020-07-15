package packets

import (
	"bytes"
	"fmt"
	"testing"
)

func TestDisconnectPacket_Read(t *testing.T) {
	var d DisconnectPacket
	disconnect := bytes.NewBuffer([]byte{0})

	err := d.Read(disconnect)
	if err != nil {
		t.Errorf("could not read %s packet: %v\n", d.name, err)
	}
	fmt.Printf("read %s packet: %+v\n", d.name, d)
}

func TestDisconnectPacket_Write(t *testing.T) {
	var buf bytes.Buffer
	var d DisconnectPacket

	d.CreatePacket()
	err := d.Write(&buf)
	if err != nil {
		t.Errorf("could not write %s packet: %v\n", d.name, err)
	}
	fmt.Printf("write %s packet: %+v\n", d.name, d)
}

func TestDisconnectPacket_Compare(t *testing.T) {
	var dRead, dDefault DisconnectPacket
	disconnect := bytes.NewBuffer([]byte{0})

	// correct packet
	err := dRead.Read(disconnect)
	if err != nil {
		t.Errorf("could not read %s packet: %v\n", dRead.name, err)
	}
	fmt.Printf("read %s packet: %+v\n", dRead.name, dRead)

	// create packet
	dDefault.CreatePacket()
	fmt.Printf("default %s packet: %+v\n", dDefault.name, dDefault)

	if dRead != dDefault {
		t.Errorf("%s packets don't match\n", dRead.name)
	}
}
