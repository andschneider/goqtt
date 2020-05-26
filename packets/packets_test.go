package packets

import (
	"bytes"
	"testing"
)

const (
	testTopic = "test/topic"
)

func TestFixedHeader(t *testing.T) {
	var expected bytes.Buffer
	expected.Write([]byte{16, 0})

	fh := FixedHeader{MessageType: "CONNECT"}
	h := fh.WriteHeader()
	if !bytes.Equal(expected.Bytes(), h.Bytes()) {
		t.Errorf("Expected %b, got %b\n", expected, h)
	}
}
