package packets

import (
	"bytes"
	"fmt"
	"testing"
)

func TestDisconnectPacket_ReadDisconnectPacket(t *testing.T) {
	disconnect := bytes.NewBuffer([]byte{0})

	d := DisconnectPacket{}
	err := d.ReadDisconnectPacket(disconnect)
	if err != nil {
		t.Errorf("could not read %s packet: %v\n", d.name, err)
	}
	fmt.Printf("%s packet: %s\n", d.name, &d)
}

func TestDisconnectPacket_Write(t *testing.T) {
	disconnect := bytes.NewBuffer([]byte{0})

	// correct packet
	d := DisconnectPacket{}
	err := d.ReadDisconnectPacket(disconnect)
	if err != nil {
		t.Errorf("could not read %s packet: %v\n", d.name, err)
	}
	fmt.Printf("%s packet: %s\n", d.name, &d)

	// create packet
	var buf bytes.Buffer
	dp := CreateDisconnectPacket()
	err = dp.Write(&buf)
	if err != nil {
		t.Errorf("could not write %s packet: %v", dp.name, err)
	}
	fmt.Printf("%s packet: %s\n", dp.name, &dp)

	if dp != d {
		t.Errorf("%s packets don't match\n", dp.name)
	}
}
