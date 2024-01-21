package pixel

import (
	"encoding/binary"
	"errors"
	"image"
	"image/color"
	"image/draw"

	"github.com/millken/x/pixel/pixelcolor"
)

type Format int

const (
	UnknownFormat Format = iota
	MHMSBFormat
	MVLSBFormat
	RGB332Format
	RGB565Format
	RGB888Format
	RGBA4444Format
	RGBA5551Format
)

func (f Format) ColorModel() color.Model {
	switch f {
	case MHMSBFormat, MVLSBFormat:
		return pixelcolor.BitModel
	case RGB332Format:
		return pixelcolor.RGB332Model
	case RGB565Format:
		return pixelcolor.RGB565Model
	case RGB888Format:
		return pixelcolor.RGB888Model
	case RGBA4444Format:
		return pixelcolor.RGBA4444Model
	case RGBA5551Format:
		return pixelcolor.RGBA5551Model
	default:
		return color.RGBAModel
	}
}

type Image interface {
	draw.Image
}

type Bitmap struct {
	Rect   image.Rectangle
	Pix    []byte
	Stride int
	Format Format
}

func NewBitmap(w, h int) *Bitmap {
	var (
		area = w * h
		pix  = area >> 3
	)
	if pix<<3 != area {
		pix++
	}
	return &Bitmap{
		Rect:   image.Rectangle{Max: image.Point{X: w, Y: h}},
		Pix:    make([]byte, pix),
		Stride: w,
		Format: MVLSBFormat,
	}
}

func (b *Bitmap) ColorModel() color.Model {
	return b.Format.ColorModel()
}

func (b *Bitmap) Bounds() image.Rectangle {
	return b.Rect
}

func (b *Bitmap) At(x, y int) color.Color {
	if x < 0 || y < 0 || x >= b.Rect.Max.X || y >= b.Rect.Max.Y {
		return pixelcolor.Off
	}
	offset, mask := b.PixOffset(x, y)
	if offset >= len(b.Pix) {
		return pixelcolor.Off
	}
	return pixelcolor.Bit(b.Pix[offset]&byte(mask) != 0)
}

func (b *Bitmap) PixOffset(x, y int) (int, uint) {
	if b.Format == MVLSBFormat {
		offset := (y>>3)*b.Stride + x
		bit := uint(y & 7)
		return offset, 1 << bit
	}
	offset := (y*b.Stride + x) >> 3
	bit := uint(7 - (y & 7))
	return offset, 1 << bit
}

func (b *Bitmap) Set(x, y int, c color.Color) {
	b.SetBit(x, y, pixelcolor.ToBit(c))
}

func (b *Bitmap) SetBit(x, y int, c pixelcolor.Bit) {
	offset, mask := b.PixOffset(x, y)
	if offset < 0 || offset >= len(b.Pix) {
		return
	}
	if c {
		b.Pix[offset] |= byte(mask)
	} else {
		b.Pix[offset] &= ^byte(mask)
	}
}

type RGB332 struct {
	Rect   image.Rectangle
	Pix    []byte
	Stride int
}

func (b *RGB332) Bounds() image.Rectangle {
	return b.Rect
}

func (b *RGB332) At(x, y int) color.Color {
	if !(image.Point{X: x, Y: y}).In(b.Rect) {
		return pixelcolor.RGB332(0)
	}
	return pixelcolor.RGB332(b.Pix[y*b.Stride+x])
}

func (b *RGB332) Set(x, y int, c color.Color) {
	if !(image.Point{X: x, Y: y}).In(b.Rect) {
		return
	}
	b.Pix[y*b.Stride+x] = byte(pixelcolor.ToRGB332(c))
}

func (RGB332) ColorModel() color.Model {
	return pixelcolor.RGB332Model
}

func NewRGB332(w, h int) *RGB332 {
	return &RGB332{
		Rect:   image.Rectangle{Max: image.Point{X: w, Y: h}},
		Stride: w,
		Pix:    make([]byte, w*h),
	}
}

type RGB565 struct {
	Rect   image.Rectangle
	Pix    []byte
	Stride int
}

func (b *RGB565) Bounds() image.Rectangle {
	return b.Rect
}

func (b *RGB565) At(x, y int) color.Color {
	if !(image.Point{X: x, Y: y}).In(b.Rect) {
		return pixelcolor.RGB565(0)
	}
	return pixelcolor.RGB565(binary.BigEndian.Uint16(b.Pix[b.OffsetOf(x, y):]))
}

func (b *RGB565) Set(x, y int, c color.Color) {
	if !(image.Point{X: x, Y: y}).In(b.Rect) {
		return
	}
	binary.BigEndian.PutUint16(b.Pix[b.OffsetOf(x, y):], uint16(pixelcolor.ToRGB565(c)))
}

func (b *RGB565) OffsetOf(x, y int) (offset int) {
	return y*b.Stride + x*2
}

func (RGB565) ColorModel() color.Model {
	return pixelcolor.RGB565Model
}

func NewRGB565(w, h int) *RGB565 {
	return &RGB565{
		Rect:   image.Rectangle{Max: image.Point{X: w, Y: h}},
		Pix:    make([]byte, w*h*2),
		Stride: w * 2,
	}
}

type RGB888 struct {
	Rect   image.Rectangle
	Pix    []byte
	Stride int
}

func (b *RGB888) Bounds() image.Rectangle {
	return b.Rect
}

func (b *RGB888) At(x, y int) color.Color {
	if !(image.Point{X: x, Y: y}).In(b.Rect) {
		return color.Black
	}
	v := b.Pix[b.OffsetOf(x, y):]
	return pixelcolor.RGB888{v[0], v[1], v[2]}
}

func (b *RGB888) Set(x, y int, c color.Color) {
	if !(image.Point{X: x, Y: y}).In(b.Rect) {
		return
	}
	v := pixelcolor.ToRGB888(c)
	copy(b.Pix[b.OffsetOf(x, y):], []byte{v.R, v.G, v.B})
}

func (b *RGB888) OffsetOf(x, y int) (offset int) {
	return y*b.Stride + x*3
}

func (b *RGB888) ColorModel() color.Model {
	return pixelcolor.RGB888Model
}

func NewRGB888(w, h int) *RGB888 {
	return &RGB888{
		Rect:   image.Rectangle{Max: image.Point{X: w, Y: h}},
		Pix:    make([]byte, w*h*3),
		Stride: w * 3,
	}
}

type RGBA4444 struct {
	Rect   image.Rectangle
	Pix    []byte
	Stride int
}

func (b *RGBA4444) Bounds() image.Rectangle {
	return b.Rect
}

func (b *RGBA4444) At(x, y int) color.Color {
	if !(image.Point{X: x, Y: y}).In(b.Rect) {
		return color.Black
	}
	return pixelcolor.RGBA4444(binary.BigEndian.Uint16(b.Pix[b.OffsetOf(x, y):]))
}

func (b *RGBA4444) Set(x, y int, c color.Color) {
	if !(image.Point{X: x, Y: y}).In(b.Rect) {
		return
	}
	binary.BigEndian.PutUint16(b.Pix[b.OffsetOf(x, y):], uint16(pixelcolor.ToRGBA4444(c)))
}

func (b *RGBA4444) OffsetOf(x, y int) (offset int) {
	return y*b.Stride + x*3
}

func (b *RGBA4444) ColorModel() color.Model {
	return pixelcolor.RGBA4444Model
}

func NewRGBA4444(w, h int) *RGBA4444 {
	return &RGBA4444{
		Rect:   image.Rectangle{Max: image.Point{X: w, Y: h}},
		Pix:    make([]byte, w*h*2),
		Stride: w * 2,
	}
}

type RGBA5551 struct {
	Rect   image.Rectangle
	Pix    []byte
	Stride int
}

func (b *RGBA5551) Bounds() image.Rectangle {
	return b.Rect
}

func (b *RGBA5551) At(x, y int) color.Color {
	if !(image.Point{X: x, Y: y}).In(b.Rect) {
		return color.Black
	}
	return pixelcolor.RGBA5551(binary.BigEndian.Uint16(b.Pix[b.OffsetOf(x, y):]))
}

func (b *RGBA5551) Set(x, y int, c color.Color) {
	if !(image.Point{X: x, Y: y}).In(b.Rect) {
		return
	}
	binary.BigEndian.PutUint16(b.Pix[b.OffsetOf(x, y):], uint16(pixelcolor.ToRGBA5551(c)))
}

func (b *RGBA5551) OffsetOf(x, y int) (offset int) {
	return y*b.Stride + x*3
}

func (b *RGBA5551) ColorModel() color.Model {
	return pixelcolor.RGBA5551Model
}

func NewRGBA5551(w, h int) *RGBA5551 {
	return &RGBA5551{
		Rect:   image.Rectangle{Max: image.Point{X: w, Y: h}},
		Pix:    make([]byte, w*h*2),
		Stride: w * 2,
	}
}

func New(width, height int, format Format) (Image, error) {
	switch format {
	case MHMSBFormat, MVLSBFormat:
		b := NewBitmap(width, height)
		b.Format = format
		return b, nil
	case RGB332Format:
		return NewRGB332(width, height), nil
	case RGB565Format:
		return NewRGB565(width, height), nil
	case RGB888Format:
		return NewRGB888(width, height), nil
	case RGBA4444Format:
		return NewRGBA4444(width, height), nil
	case RGBA5551Format:
		return NewRGBA5551(width, height), nil
	default:
		return nil, errors.New("framebuffer: invalid format")
	}
}

func Must(width, height int, format Format) Image {
	b, err := New(width, height, format)
	if err != nil {
		panic(err)
	}
	return b
}
