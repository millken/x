package pixelfont

import (
	"image"
	"image/draw"
	"io"
	"io/ioutil"
	"os"

	"github.com/millken/x/pixel"
	"github.com/millken/x/pixel/pixelcolor"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// Builtin fonts.
var (
	GLCD5x8 = NewFont(5, 8, fontGLCD5x8)
)

// Font can render glyph images.
type Font interface {
	// Glyph returns the image mask for the requested glyph.
	Glyph(r rune) image.Image

	// Bounds is the bounding box that fits any glyph in the font.
	Bounds() image.Rectangle
}

type cachedFont struct {
	Font
	cache map[rune]image.Image
}

func (c *cachedFont) Glyph(r rune) image.Image {
	if _, ok := c.cache[r]; !ok {
		c.cache[r] = c.Font.Glyph(r)
	}
	return c.cache[r]
}

// NewFontCache is a helper that caches all rendered glyphs in memory.
func NewFontCache(font Font) Font {
	return &cachedFont{
		Font:  font,
		cache: make(map[rune]image.Image),
	}
}

// TexturedFont applies a texture to font bitmaps in the output glyph.
type TexturedFont struct {
	Font
	Texture image.Image
}

func (f TexturedFont) Glyph(r rune) image.Image {
	var (
		b     = f.Bounds()
		glyph = image.NewRGBA(b)
		mask  = f.Font.Glyph(r)
	)
	if mask != nil && f.Texture != nil {
		draw.DrawMask(glyph, b, f.Texture, image.Point{}, mask, image.Point{}, draw.Src)
	}
	return glyph
}

type bitmapFont struct {
	rect   image.Rectangle
	stride int
	pix    []byte
}

func (f *bitmapFont) Glyph(r rune) image.Image {
	o := int(r) * f.stride
	if o < 0 || o >= len(f.pix) {
		return nil
	}
	b := pixel.NewBitmap(f.rect.Max.X, f.rect.Max.Y)
	for x := 0; x < f.rect.Max.X; x++ {
		v := f.pix[o+x]
		for y := 0; y < f.rect.Max.Y; y++ {
			b.Set(x, y, pixelcolor.Bit((v>>y)&0x01 == 0x01))
		}
	}
	return b
}

func (f *bitmapFont) Bounds() image.Rectangle {
	return f.rect
}

func NewFont(w, h int, pix []byte) Font {
	return &bitmapFont{
		rect:   image.Rectangle{Max: image.Point{X: w, Y: h}},
		stride: w,
		pix:    pix,
	}
}

func LoadFont(name string) (Font, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	var dim [2]byte
	if _, err = io.ReadFull(f, dim[:]); err != nil {
		_ = f.Close()
		return nil, err
	}

	var pix []byte
	if pix, err = ioutil.ReadAll(f); err != nil {
		_ = f.Close()
		return nil, err
	}
	if err = f.Close(); err != nil {
		return nil, err
	}

	return NewFont(int(dim[0]), int(dim[1]), pix), nil
}

type ttfFont struct {
	face font.Face
	rect image.Rectangle
}

func (f ttfFont) Glyph(r rune) image.Image {
	/*
		var (
			size  = measureRune(f.face, r)
			glyph = NewBitmap(size.X, size.Y)
			//center = fixed.P((f.rect.Max.X-size.X)/2, (f.rect.Max.Y-size.Y)/2+size.Y+1)
			center = fixed.P((f.rect.Max.X-size.X)/2, size.Y)
		)
		log.Printf("glyph %c (%d) is %dx%d", r, r, size.X, size.Y)
		dr, mask, maskPoint, _, ok := f.face.Glyph(center, r)
		if !ok {
			return nil
		}
	*/
	dr, mask, maskPoint, _, ok := f.face.Glyph(fixed.Point26_6{}, r)
	if !ok {
		return nil
	}
	glyph := pixel.NewBitmap(dr.Dx(), dr.Dy())
	// log.Printf("glyph: %s, dr: %s, mask: %s, mask point: %s", glyph.Bounds(), dr, mask.Bounds(), maskPoint)
	draw.DrawMask(glyph, glyph.Bounds(), image.White, image.Point{}, mask, maskPoint, draw.Over)
	return glyph
}

func (f ttfFont) Bounds() image.Rectangle {
	return f.rect
}

func NewTTFFont(ttf []byte, size int) (Font, error) {
	f, err := truetype.Parse(ttf)
	if err != nil {
		return nil, err
	}
	face := truetype.NewFace(f, &truetype.Options{
		Size:    float64(size),
		DPI:     72, // At 72 DPI, px == pt
		Hinting: font.HintingNone,
	})
	return &ttfFont{
		face: face,
		rect: measureFont(face, size),
	}, nil
}

func LoadTTFFont(name string, size int) (Font, error) {
	ttf, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return NewTTFFont(ttf, size)
}

func measureFont(face font.Face, height int) (b image.Rectangle) {
	b.Max.Y = height
	const testGlyphs = "xXmMjJ0"
	for _, r := range testGlyphs {
		size := measureRune(face, r)
		b.Max.X = max(b.Max.X, size.X)
		b.Max.Y = max(b.Max.Y, size.Y)
	}
	return
}

func measureRune(face font.Face, r rune) image.Point {
	bounds, advance, ok := face.GlyphBounds(r)
	if !ok {
		return image.Point{}
	}
	return image.Pt(advance.Ceil(), (bounds.Max.Y - bounds.Min.Y).Ceil())
}
