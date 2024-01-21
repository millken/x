package pixeleffect

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	_ "image/png" // PNG codec

	"github.com/millken/x/pixel"
)

func TestDitherStep(t *testing.T) {
	im := testLoadImage(t, filepath.Join("testdata", "gopher.png"))
	for i := 0; i < len(ditherMasks); i++ {
		DitherStep(im, i)
		if os.Getenv("DEBUG_DITHERSTEP") != "" {
			o, err := os.Create(filepath.Join(os.TempDir(), fmt.Sprintf("dither-step-%d.png", i)))
			if err != nil {
				t.Fatal(err)
			}
			if err = png.Encode(o, im); err != nil {
				t.Fatal(err)
			}
			if err = o.Close(); err != nil {
				t.Fatal(err)
			}
			t.Log("saved to", o.Name())
		}
	}
}

func TestFadeOutVertical(t *testing.T) {
	var (
		im     = testLoadImage(t, filepath.Join("testdata", "gopher.png"))
		bounds = im.Bounds()
		width  = bounds.Dx()
		height = bounds.Dy()
		steps  = width / 2
	)
	if width&1 == 1 {
		steps++
	}
	for i := 0; i < im.Bounds().Dy()/2; i++ {
		offset := i * 2
		pixel.VLine(im, image.Point{X: offset}, height, color.Transparent)
		pixel.VLine(im, image.Point{X: height - offset - 1}, height, color.Transparent)
	}
}

func testLoadImage(t *testing.T, name string) pixel.Image {
	f, err := os.Open(name)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	i, _, err := image.Decode(f)
	if err != nil {
		t.Fatal("error decoding", name+":", err)
	}
	return toRGBA(i)
}

func toRGBA(i image.Image) *image.RGBA {
	if o, ok := i.(*image.RGBA); ok {
		return o
	}
	o := image.NewRGBA(i.Bounds())
	draw.Draw(o, o.Bounds(), i, image.Point{}, draw.Src)
	return o
}
