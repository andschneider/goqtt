package packets

import (
	"bytes"
	"fmt"
	"testing"
)

func TestConnectPacket(t *testing.T) {
	var buf bytes.Buffer
	cp := CreateConnectPacket()
	err := cp.Write(&buf, true)
	if err != nil {
		t.Errorf("could not write CONNECT packet: %v", err)
	}
	fmt.Printf("connect packet: %s\n", &cp)
}
