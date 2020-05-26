package packets

import (
	"bytes"
	"fmt"
	"testing"
)

func TestSubackPacket(t *testing.T) {
	// TODO this packet isn't totally right
	suback := bytes.NewBuffer([]byte{144, 3, 0, 1, 0})
	fmt.Printf("%08b\n", suback)

	sa := SubackPacket{}
	err := sa.ReadSubackPacket(suback)
	if err != nil {
		t.Errorf("could not read suback packet: %v\n", err)
	}
	fmt.Println(sa)
}
