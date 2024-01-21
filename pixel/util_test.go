package pixel

import (
	"bytes"
	"fmt"
	"testing"
)

func TestMemset(t *testing.T) {
	tests := []byte{0x00, 0x2a, 0xff}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%#02x", test), func(it *testing.T) {
			b := make([]byte, 512)
			memset(b, test)
			for i, v := range b {
				if v != test {
					it.Fatalf("expected byte %d to be %#02x, got %#02x", i, test, v)
				}
			}
		})
	}
}

func TestMemsetSlice(t *testing.T) {
	tests := []struct {
		Test []byte
		Size int
	}{
		{[]byte{0x00}, 128},
		{[]byte{0xff, 0xff}, 512},
		{[]byte{0x2a, 0x42, 0xaa, 0xff}, 1024},
	}
	for _, test := range tests {
		t.Run("", func(it *testing.T) {
			b := make([]byte, test.Size)
			memsetSlice(b, test.Test)
			for i := 0; i < test.Size; i += len(test.Test) {
				if !bytes.Equal(b[i:i+len(test.Test)], test.Test) {
					it.Fatalf("expected bytes %d-%d to be %x, got %x", i, i+len(test.Test), test.Test, b[i:i+len(test.Test)])
				}
			}
		})
	}
}
