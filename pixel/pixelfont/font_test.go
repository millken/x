package pixelfont

import (
	"image"
	"image/draw"
	"path/filepath"
	"testing"

	"github.com/millken/x/pixel"
	"github.com/millken/x/pixel/pixelcolor"
)

func TestBitmapFont(t *testing.T) {
	testFont(t, GLCD5x8, 5, 8)
}

func TestLoadTTFFont(t *testing.T) {
	tests := []struct {
		Name string
		Size int
	}{
		{"16bit.ttf", 16},
		{"pixelmix.ttf", 10},
	}
	for _, test := range tests {
		t.Run(test.Name, func(it *testing.T) {
			f, err := LoadTTFFont(filepath.Join("testdata", test.Name), test.Size)
			if err != nil {
				it.Fatal(err)
			}
			testFont(it, f, test.Size, test.Size)
		})
	}
}

func testFont(t *testing.T, f Font, w, h int) {
	t.Helper()
	var (
		bounds = image.Rect(0, 0, (f.Bounds().Dx()+1)*4, h)
		test   = image.NewRGBA(bounds)
		offset image.Point
	)
	draw.Draw(test, bounds, image.Black, image.Point{}, draw.Src)
	for _, r := range "Test" {
		glyph := f.Glyph(r)
		if glyph == nil {
			t.Errorf("glyph %c returned nil", r)
			continue
		}
		draw.Draw(test, bounds.Add(offset), glyph, image.Point{}, draw.Over)
		offset.X += f.Bounds().Dx() + 1
	}
	logBitmap(t, test)
}

func logBitmap(t *testing.T, b pixel.Image) {
	t.Helper()
	var (
		bounds = b.Bounds()
		row    = make([]rune, bounds.Dx())
	)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if pixelcolor.BitModel.Convert(b.At(x, y)).(pixelcolor.Bit) {
				row[x] = '#'
			} else {
				row[x] = ' '
			}
		}
		t.Logf("row %02d: %s", y+1, string(row))
	}
}
