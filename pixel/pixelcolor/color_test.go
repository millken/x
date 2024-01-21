package pixelcolor

import (
	"image/color"
	"math/bits"
	"testing"
)

var (
	testBlack     = color.RGBA{0x00, 0x00, 0x00, 0x00}
	testWhite     = color.RGBA{0xff, 0xff, 0xff, 0x00}
	testAmber     = color.RGBA{0xff, 0x7f, 0x00, 0x00}
	testIndigo    = color.RGBA{0x50, 0x00, 0x82, 0x00}
	testTurquoise = color.RGBA{0x00, 0xCE, 0xD1, 0x00}
	testGray25    = color.RGBA{0x3f, 0x3f, 0x3f, 0x3f}
	testGray50    = color.RGBA{0x7f, 0x7f, 0x7f, 0x7f}
	testGray75    = color.RGBA{0xbf, 0xbf, 0xbf, 0xbf}
	testColorSink color.Color
	testColors    = []struct {
		Name  string
		Color color.Color
	}{
		{"black", testBlack},         // black
		{"white", testWhite},         // white
		{"amber", testAmber},         // amber
		{"indigo", testIndigo},       // indigo
		{"turquoise", testTurquoise}, // turquoise
		{"gray25", testGray25},       // gray 25%
		{"gray50", testGray50},       // gray 50%
		{"gray75", testGray75},       // gray 75%
	}
)

func TestRGB332(t *testing.T) {
	tests := []struct {
		Name string
		Test RGB332
		Want color.RGBA
	}{
		{"black", 0b000_000_00, testBlack},
		{"white", 0b111_111_11, testWhite},
		{"amber", 0b111_011_00, testAmber},
		{"indigo", 0b010_000_10, testIndigo},
		{"turquoise", 0b000_110_11, testTurquoise},
		{"gray25", 0b001_001_00, testGray25},
		{"gray50", 0b011_011_01, testGray50},
		{"gray75", 0b101_101_10, testGray75},
	}
	for _, test := range tests {
		t.Run(test.Name, func(it *testing.T) {
			testColorSink = color.RGBAModel.Convert(test.Test)
			tr, tg, tb, _ := testColorSink.RGBA()
			it.Logf("RGB332(%#08b) -> RGB(%02x%02x%02x), want RGB(%02x%02x%02x)", test.Test,
				tr&0xff, tg&0xff, tb&0xff, test.Want.R, test.Want.G, test.Want.B)
		})
	}
}

func TestRGB332Model(t *testing.T) {
	want := map[string]RGB332{
		"black":     0b000_000_00,
		"white":     0b111_111_11,
		"amber":     0b111_011_00,
		"indigo":    0b010_000_10,
		"turquoise": 0b000_110_11,
		"gray25":    0b001_001_00,
		"gray50":    0b011_011_01,
		"gray75":    0b101_101_10,
	}
	for _, test := range testColors {
		t.Run(test.Name, func(it *testing.T) {
			v := RGB332Model.Convert(test.Color).(RGB332)
			if v != want[test.Name] {
				it.Fatalf("expected %q (%+v) to return %#08b, got %#08b", test.Name, test.Color, want[test.Name], v)
			}
		})
	}
}

func TestRGB565(t *testing.T) {
	tests := []struct {
		Name string
		Test RGB565
		Want color.RGBA
	}{
		{"black", 0b00000_000000_00000, testBlack},
		{"white", 0b11111_111111_11111, testWhite},
		{"amber", 0b11111_011111_00000, testAmber},
		{"indigo", 0b01010_000000_10000, testIndigo},
		{"turquoise", 0b00000_110011_11010, testTurquoise},
		{"gray25", 0b00111_001111_00111, testGray25},
		{"gray50", 0b01111_011111_01111, testGray50},
		{"gray75", 0b10111_101111_10111, testGray75},
	}
	for _, test := range tests {
		t.Run(test.Name, func(it *testing.T) {
			v := color.RGBAModel.Convert(test.Test).(color.RGBA)
			testColorBitErrors(it, test.Want, v, 3, false)
		})
	}
}

func TestRGB565Model(t *testing.T) {
	want := map[string]RGB565{
		"black":     0b00000_000000_00000,
		"white":     0b11111_111111_11111,
		"amber":     0b11111_011111_00000,
		"indigo":    0b01010_000000_10000,
		"turquoise": 0b00000_110011_11010,
		"gray25":    0b00111_001111_00111,
		"gray50":    0b01111_011111_01111,
		"gray75":    0b10111_101111_10111,
	}
	for _, test := range testColors {
		t.Run(test.Name, func(it *testing.T) {
			v := RGB565Model.Convert(test.Color).(RGB565)
			if v != want[test.Name] {
				it.Fatalf("expected %q (%+v) to return %#016b, got %#016b", test.Name, test.Color, want[test.Name], v)
			}
		})
	}
}

func TestRGB888(t *testing.T) {
	tests := []struct {
		Name string
		Test RGB888
		Want color.RGBA
	}{
		{"black", RGB888{0x00, 0x00, 0x00}, testBlack},
		{"white", RGB888{0xff, 0xff, 0xff}, testWhite},
		{"amber", RGB888{0xff, 0x7f, 0x00}, testAmber},
		{"indigo", RGB888{0x50, 0x00, 0x82}, testIndigo},
		{"turquoise", RGB888{0x00, 0xCE, 0xD1}, testTurquoise},
		{"gray25", RGB888{0x3f, 0x3f, 0x3f}, testGray25},
		{"gray50", RGB888{0x7f, 0x7f, 0x7f}, testGray50},
		{"gray75", RGB888{0xbf, 0xbf, 0xbf}, testGray75},
	}
	for _, test := range tests {
		t.Run(test.Name, func(it *testing.T) {
			v := color.RGBAModel.Convert(test.Test).(color.RGBA)
			testColorBitErrors(it, test.Want, v, 2, false)
		})
	}
}

func TestRGB888Model(t *testing.T) {
	want := map[string]RGB888{
		"black":     {0x00, 0x00, 0x00}, // black
		"white":     {0xff, 0xff, 0xff}, // white
		"amber":     {0xff, 0x7f, 0x00}, // amber
		"indigo":    {0x50, 0x00, 0x82}, // indigo
		"turquoise": {0x00, 0xCE, 0xD1}, // turquoise
		"gray25":    {0x3f, 0x3f, 0x3f}, // gray 25%
		"gray50":    {0x7f, 0x7f, 0x7f}, // gray 50%
		"gray75":    {0xbf, 0xbf, 0xbf}, // gray 75%
	}
	for _, test := range testColors {
		t.Run(test.Name, func(it *testing.T) {
			v := RGB888Model.Convert(test.Color).(RGB888)
			if v != want[test.Name] {
				it.Fatalf("expected %q (%+v) to return %+v, got %+v", test.Name, test.Color, want[test.Name], v)
			}
		})
	}
}

func TestRGBA4444(t *testing.T) {
	tests := []struct {
		Name string
		Test RGBA4444
		Want color.RGBA
	}{
		{"black", 0b0000_0000_0000_0000, testBlack},
		{"white", 0b1111_1111_1111_0000, testWhite},
		{"amber", 0b1111_0111_0000_0000, testAmber},
		{"indigo", 0b0101_0000_1000_0000, testIndigo},
		{"turquoise", 0b0000_1100_1101_0000, testTurquoise},
		{"gray25", 0b0011_0011_0011_0011, testGray25},
		{"gray50", 0b0111_0111_0111_0111, testGray50},
		{"gray75", 0b1011_1011_1011_1011, testGray75},
	}
	for _, test := range tests {
		t.Run(test.Name, func(it *testing.T) {
			v := color.RGBAModel.Convert(test.Test).(color.RGBA)
			testColorBitErrors(it, test.Want, v, 2, true)
		})
	}
}

func TestRGBA444Model(t *testing.T) {
	want := map[string]RGBA4444{
		"black":     0b0000_0000_0000_0000, // black
		"white":     0b1111_1111_1111_0000, // white
		"amber":     0b1111_0111_0000_0000, // amber
		"indigo":    0b0101_0000_1000_0000, // indigo
		"turquoise": 0b0000_1100_1101_0000, // turquoise
		"gray25":    0b0011_0011_0011_0011, // gray 50%
		"gray50":    0b0111_0111_0111_0111, // gray 50%
		"gray75":    0b1011_1011_1011_1011, // gray 50%
	}
	for _, test := range testColors {
		t.Run(test.Name, func(it *testing.T) {
			v := RGBA4444Model.Convert(test.Color).(RGBA4444)
			if v != want[test.Name] {
				it.Fatalf("expected %q (%+v) to return %#016b, got %#016b", test.Name, test.Color, want[test.Name], v)
			}
		})
	}
}

func TestRGBA5551(t *testing.T) {
	tests := []struct {
		Name  string
		Test  RGBA5551
		Want  color.RGBA // A not used
		Alpha bool
	}{
		{"black", 0b00000_00000_00000_0, testBlack, false},
		{"white", 0b11111_11111_11111_0, testWhite, false},
		{"amber", 0b11111_01111_00000_0, testAmber, false},
		{"indigo", 0b01010_00000_10000_0, testIndigo, false},
		{"turquoise", 0b00000_11001_11010_0, testTurquoise, false},
		{"gray25", 0b00111_00111_00111_1, testGray25, true},
		{"gray50", 0b01111_01111_01111_1, testGray50, true},
		{"gray75", 0b10111_10111_10111_1, testGray75, true},
	}
	for _, test := range tests {
		t.Run(test.Name, func(it *testing.T) {
			v := color.RGBAModel.Convert(test.Test).(color.RGBA)
			testColorBitErrors(it, test.Want, v, 3, false)
			if v.A > 0 != test.Alpha {
				it.Errorf("expected alpha %t, got %t: %+v", test.Alpha, v.A > 0, v)
			}
		})
	}
}

func TestRGBA5551Model(t *testing.T) {
	want := map[string]RGBA5551{
		"black":     0b00000_00000_00000_0, // black
		"white":     0b11111_11111_11111_0, // white
		"amber":     0b11111_01111_00000_0, // amber
		"indigo":    0b01010_00000_10000_0, // indigo
		"turquoise": 0b00000_11001_11010_0, // turquoise
		"gray25":    0b00111_00111_00111_1, // gray 25%
		"gray50":    0b01111_01111_01111_1, // gray 50%
		"gray75":    0b10111_10111_10111_1, // gray 75%
	}
	for _, test := range testColors {
		t.Run(test.Name, func(it *testing.T) {
			v := RGBA5551Model.Convert(test.Color).(RGBA5551)
			if v != want[test.Name] {
				it.Fatalf("expected %q (%+v) to return %#016b, got %#016b", test.Name, test.Color, want[test.Name], v)
			}
		})
	}
}

func testColorBitErrors(t *testing.T, want color.RGBA, v color.RGBA, errors int, alpha bool) {
	t.Helper()
	if n := bits.OnesCount8(v.R ^ want.R); n > errors {
		t.Errorf("R channel has %d bit errors: want %#02x (%#08b), got %#02x (%#08b)", n, want.R, want.R, v.R, v.R)
	}
	if n := bits.OnesCount8(v.G ^ want.G); n > errors {
		t.Errorf("G channel has %d bit errors: want %#02x (%#08b), got %#02x (%#08b)", n, want.G, want.G, v.G, v.G)
	}
	if n := bits.OnesCount8(v.B ^ want.B); n > errors {
		t.Errorf("B channel has %d bit errors: want %#02x (%#08b), got %#02x (%#08b)", n, want.B, want.B, v.B, v.B)
	}
	if alpha {
		if n := bits.OnesCount8(v.A ^ want.A); n > errors {
			t.Errorf("A channel has %d bit errors: want %#02x (%#08b), got %#02x (%#08b)", n, want.A, want.A, v.A, v.A)
		}
	}
}
