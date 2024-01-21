package pixel

import (
	"fmt"
	"testing"
)

func TestNewBitBuffer(t *testing.T) {
	tests := []struct {
		W, H       int
		WantStride int
		WantPixLen int
	}{
		{1, 1, 1, 1},
		{3, 3, 3, 2},
		{8, 8, 8, 8},
		{128, 32, 128, 512},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%dx%d", test.W, test.H), func(it *testing.T) {
			v := NewBitmap(test.W, test.H)
			if v == nil {
				it.Fatal("NewBitmap returned nil")
			}
			if v.Stride != test.WantStride {
				it.Errorf("expected stride %d, got %d", test.WantStride, v.Stride)
			}
			if l := len(v.Pix); l != test.WantPixLen {
				it.Errorf("expected pix len %d, got %d", test.WantPixLen, l)
			}
		})
	}
}

func TestNewRGB565Buffer(t *testing.T) {
	tests := []struct {
		W, H       int
		WantStride int
		WantPixLen int
	}{
		{8, 8, 16, 128},
		{128, 32, 256, 8192},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%dx%d", test.W, test.H), func(it *testing.T) {
			v := NewRGB565(test.W, test.H)
			if v == nil {
				it.Fatal("NewRGB565 returned nil")
			}
			if v.Stride != test.WantStride {
				it.Errorf("expected stride %d, got %d", test.WantStride, v.Stride)
			}
			if l := len(v.Pix); l != test.WantPixLen {
				it.Errorf("expected pix len %d, got %d", test.WantPixLen, l)
			}
		})
	}
}

func TestNewRGB888Buffer(t *testing.T) {
	tests := []struct {
		W, H       int
		WantStride int
		WantPixLen int
	}{
		{8, 8, 24, 192},
		{128, 32, 384, 12288},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%dx%d", test.W, test.H), func(it *testing.T) {
			v := NewRGB888(test.W, test.H)
			if v == nil {
				it.Fatal("NewRGB888 returned nil")
			}
			if v.Stride != test.WantStride {
				it.Errorf("expected stride %d, got %d", test.WantStride, v.Stride)
			}
			if l := len(v.Pix); l != test.WantPixLen {
				it.Errorf("expected pix len %d, got %d", test.WantPixLen, l)
			}
		})
	}
}
