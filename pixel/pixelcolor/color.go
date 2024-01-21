//go:generate go run color_table_gen.go

package pixelcolor

import (
	"image/color"
)

// Color models supported by this package.
var (
	BitModel         = color.ModelFunc(bitModel)
	RGB332Model      = color.ModelFunc(rgb332Model)
	RGB565Model      = color.ModelFunc(rgb565Model)
	RGB888Model      = color.ModelFunc(rgb888Model)
	RGBA2222Model    = color.ModelFunc(rgba2222Model)
	RGBA3323Model    = color.ModelFunc(rgba3323Model)
	RGBA4444Model    = color.ModelFunc(rgba4444Model)
	RGBA5551Model    = color.ModelFunc(rgba5551Model)
	RGBX5551Model    = color.ModelFunc(rgbx5551Model)
	RGBA5656Model    = color.ModelFunc(rgba5656Model)
	RGBA8888Model    = color.RGBAModel
	RGBA1010102Model = color.ModelFunc(rgba1010102Model)
	ARGB1555Model    = color.ModelFunc(argb1555Model)
	ARGB2222Model    = color.ModelFunc(argb2222Model)
	ARGB3332Model    = color.ModelFunc(argb3332Model)
	ARGB4444Model    = color.ModelFunc(argb4444Model)
	ARGB6565Model    = color.ModelFunc(argb6565Model)
	ARGB8888Model    = color.ModelFunc(argb8888Model)
	ARGB2101010Model = color.ModelFunc(argb2101010Model)
	XRGB1555Model    = color.ModelFunc(xrgb1555Model)
)

func bitModel(c color.Color) color.Color         { return ToBit(c) }
func rgb332Model(c color.Color) color.Color      { return ToRGB332(c) }
func rgb565Model(c color.Color) color.Color      { return ToRGB565(c) }
func rgb888Model(c color.Color) color.Color      { return ToRGB888(c) }
func rgba2222Model(c color.Color) color.Color    { return ToRGBA2222(c) }
func rgba3323Model(c color.Color) color.Color    { return ToRGBA3323(c) }
func rgba4444Model(c color.Color) color.Color    { return ToRGBA4444(c) }
func rgba5551Model(c color.Color) color.Color    { return ToRGBA5551(c) }
func rgbx5551Model(c color.Color) color.Color    { return ToRGBX5551(c) }
func rgba5656Model(c color.Color) color.Color    { return ToRGBA5656(c) }
func rgba1010102Model(c color.Color) color.Color { return ToRGBA1010102(c) }
func argb1555Model(c color.Color) color.Color    { return ToARGB1555(c) }
func argb2222Model(c color.Color) color.Color    { return ToARGB2222(c) }
func argb3332Model(c color.Color) color.Color    { return ToARGB3332(c) }
func argb4444Model(c color.Color) color.Color    { return ToARGB4444(c) }
func argb6565Model(c color.Color) color.Color    { return ToARGB6565(c) }
func argb8888Model(c color.Color) color.Color    { return ToARGB8888(c) }
func argb2101010Model(c color.Color) color.Color { return ToARGB2101010(c) }
func xrgb1555Model(c color.Color) color.Color    { return ToXRGB1555(c) }

// Bit is a 1-bit color.
type Bit bool

// Bit values.
const (
	Off Bit = false
	On  Bit = true
)

func (c Bit) RGBA() (r, g, b, a uint32) {
	if c {
		return 0xffff, 0xffff, 0xffff, 0xffff
	}
	return 0, 0, 0, 0xffff
}

func ToBit(c color.Color) Bit {
	switch c := c.(type) {
	case Bit:
		return c
	default:
		r, g, b, _ := c.RGBA()

		// These coefficients (the fractions 0.299, 0.587 and 0.114) are the same
		// as those given by the JFIF specification and used by func RGBToYCbCr in
		// ycbcr.go.
		//
		// Note that 19595 + 38470 + 7471 equals 65536.
		y := (19595*r + 38470*g + 7471*b + 1<<15) >> 24

		return y >= 0x80
	}
}

// RGB332 is a 8-bit RGB color with no alpha channel.
type RGB332 uint8 // 3-3-2 RGB

func (c RGB332) RGBA() (r, g, b, a uint32) {
	r = lut3to8[(c&0b111_000_00)>>5]
	g = lut3to8[(c&0b000_111_00)>>2]
	b = lut2to8[(c&0b000_000_11)>>0]
	a = 0xffff
	return
}

func ToRGB332(c color.Color) RGB332 {
	if c, ok := c.(RGB332); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	r = lut8to3[(r>>8)&0xff]
	g = lut8to3[(g>>8)&0xff]
	b = lut8to2[(b>>8)&0xff]
	return RGB332(r<<5 | g<<2 | b)
}

// RGB565 is a 16-bit RGB color with no alpha channel.
type RGB565 uint16 // 5-6-5 RGB

func (c RGB565) RGBA() (r, g, b, a uint32) {
	r = lut5to8[(c&0b11111_000000_00000)>>11]
	g = lut6to8[(c&0b00000_111111_00000)>>5]
	b = lut5to8[(c&0b00000_000000_11111)>>0]
	a = 0xffff
	return
}

func ToRGB565(c color.Color) RGB565 {
	r, g, b, _ := c.RGBA()
	r = lut8to5[(r>>8)&0xff]
	g = lut8to6[(g>>8)&0xff]
	b = lut8to5[(b>>8)&0xff]
	return RGB565(r<<11 | g<<5 | b)
}

// RGB888 is a 24-bit color with no alpha channel.
type RGB888 struct {
	R, G, B uint8
}

func (c RGB888) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R)
	r |= r << 8
	g = uint32(c.G)
	g |= g << 8
	b = uint32(c.B)
	b |= b << 8
	a = 0xffff
	return
}

func ToRGB888(c color.Color) RGB888 {
	r, g, b, _ := c.RGBA()
	return RGB888{byte(r), byte(g), byte(b)}
}

// RGBA2222 is a 8-bit RGB color with alpha channel.
type RGBA2222 uint8

func (c RGBA2222) RGBA() (r, g, b, a uint32) {
	r = lut2to8[uint32(c&0b11_00_00_00)>>6]
	g = lut2to8[uint32(c&0b00_11_00_00)>>4]
	b = lut2to8[uint32(c&0b00_00_11_00)>>2]
	a = lut2to8[uint32(c&0b00_00_00_11)>>0]
	return
}

func ToRGBA2222(c color.Color) RGBA2222 {
	if c, ok := c.(RGBA2222); ok {
		return c
	}
	r, g, b, a := c.RGBA()
	r = lut8to2[(r>>8)&0xff]
	g = lut8to2[(g>>8)&0xff]
	b = lut8to2[(b>>8)&0xff]
	a = lut8to2[(a>>8)&0xff]
	return RGBA2222(r<<6 | g<<4 | b<<2 | a)
}

// RGBA3323 is a 8-bit RGB color with alpha channel.
type RGBA3323 uint16

func (c RGBA3323) RGBA() (r, g, b, a uint32) {
	r = lut3to8[(c&0b111_000_00_000)>>8]
	g = lut3to8[(c&0b000_111_00_000)>>5]
	b = lut2to8[(c&0b000_000_11_000)>>3]
	a = lut3to8[(c&0b000_000_00_111)>>0]
	return
}

func ToRGBA3323(c color.Color) RGBA3323 {
	if c, ok := c.(RGBA3323); ok {
		return c
	}
	r, g, b, a := c.RGBA()
	r = lut8to3[(r>>8)&0xff]
	g = lut8to3[(g>>8)&0xff]
	b = lut8to2[(b>>8)&0xff]
	a = lut8to3[(a>>8)&0xff]
	return RGBA3323(r<<8 | g<<5 | b<<3 | a)
}

// RGBA4444 is a 16-bit RGB color with alpha channel.
type RGBA4444 uint16

func (c RGBA4444) RGBA() (r, g, b, a uint32) {
	r = lut4to8[(c&0b1111_0000_0000_0000)>>12]
	g = lut4to8[(c&0b0000_1111_0000_0000)>>8]
	b = lut4to8[(c&0b0000_0000_1111_0000)>>4]
	a = lut4to8[(c&0b0000_0000_0000_1111)>>0]
	return
}

func ToRGBA4444(c color.Color) RGBA4444 {
	if c, ok := c.(RGBA4444); ok {
		return c
	}
	r, g, b, a := c.RGBA()
	r = lut8to4[(r>>8)&0xff]
	g = lut8to4[(g>>8)&0xff]
	b = lut8to4[(b>>8)&0xff]
	a = lut8to4[(a>>8)&0xff]
	return RGBA4444(r<<12 | g<<8 | b<<4 | a)
}

// RGBA5551 is a 16-bit RGB color with alpha channel.
type RGBA5551 uint16

func (c RGBA5551) RGBA() (r, g, b, a uint32) {
	r = lut5to8[(c&0b11111_00000_00000_0)>>11]
	g = lut5to8[(c&0b00000_11111_00000_0)>>6]
	b = lut5to8[(c&0b00000_00000_11111_0)>>1]
	if (c & 0b00000_00000_00000_1) == 1 {
		a = 0xffff
	}
	return
}

func ToRGBA5551(c color.Color) RGBA5551 {
	if c, ok := c.(RGBA5551); ok {
		return c
	}
	r, g, b, a := c.RGBA()
	r = lut8to5[(r>>8)&0xff]
	g = lut8to5[(g>>8)&0xff]
	b = lut8to5[(b>>8)&0xff]
	if a > 0 {
		a = 1
	}
	return RGBA5551(r<<11 | g<<6 | b<<1 | a)
}

// RGBX5551 is a 16-bit RGB color without alpha channel.
type RGBX5551 uint16

func (c RGBX5551) RGBA() (r, g, b, a uint32) {
	r = lut5to8[(c&0b11111_00000_00000_0)>>11]
	g = lut5to8[(c&0b00000_11111_00000_0)>>6]
	b = lut5to8[(c&0b00000_00000_11111_0)>>1]
	a = 0xffff
	return
}

func ToRGBX5551(c color.Color) RGBX5551 {
	if c, ok := c.(RGBX5551); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	r = lut8to5[(r>>8)&0xff]
	g = lut8to5[(g>>8)&0xff]
	b = lut8to5[(b>>8)&0xff]
	return RGBX5551(r<<11 | g<<6 | b<<1)
}

// RGBA5656 is a 22-bit RGB color with alpha channel.
type RGBA5656 uint32

func (c RGBA5656) RGBA() (r, g, b, a uint32) {
	r = lut5to8[(c&0b11111_000000_00000_000000)>>17]
	g = lut6to8[(c&0b00000_111111_00000_000000)>>11]
	b = lut5to8[(c&0b00000_000000_11111_000000)>>6]
	a = lut6to8[(c&0b00000_000000_00000_111111)>>0]
	return
}

func ToRGBA5656(c color.Color) RGBA5656 {
	if c, ok := c.(RGBA5656); ok {
		return c
	}
	r, g, b, a := c.RGBA()
	r = lut8to5[(r>>8)&0xff]
	g = lut8to6[(g>>8)&0xff]
	b = lut8to5[(b>>8)&0xff]
	a = lut8to6[(a>>8)&0xff]
	return RGBA5656(r<<17 | g<<11 | b<<6 | a)
}

// RGBA1010102 is a 32-bit RGB color with alpha channel.
type RGBA1010102 uint32

func (c RGBA1010102) RGBA() (r, g, b, a uint32) {
	r = uint32(c>>22) & 0x03ff
	r = r<<6 | r
	g = uint32(c>>12) & 0x03ff
	g = g<<6 | g
	b = uint32(c>>2) & 0x03ff
	b = b<<6 | b
	a = lut2to8[c&0x0003]
	return
}

func ToRGBA1010102(c color.Color) RGBA1010102 {
	if c, ok := c.(RGBA1010102); ok {
		return c
	}
	r, g, b, a := c.RGBA()
	r = (r >> 6) & 0x03ff
	g = (g >> 6) & 0x03ff
	b = (b >> 6) & 0x03ff
	a = (a >> 14) & 0x0003
	return RGBA1010102(r<<22 | g<<12 | b<<2 | a)
}

// ARGB1555 is a 16-bit RGB color with alpha channel.
type ARGB1555 uint16

func (c ARGB1555) RGBA() (r, g, b, a uint32) {
	if (c&0b1_00000_00000_00000)>>15 == 1 {
		a = 0xffff
	}
	r = lut5to8[(c&0b0_11111_00000_00000)>>10]
	g = lut5to8[(c&0b0_00000_11111_00000)>>5]
	b = lut5to8[(c&0b0_00000_00000_11111)>>0]
	return
}

func ToARGB1555(c color.Color) ARGB1555 {
	if c, ok := c.(ARGB1555); ok {
		return c
	}
	r, g, b, a := c.RGBA()
	r = lut8to5[(r>>8)&0xff]
	g = lut8to5[(g>>8)&0xff]
	b = lut8to5[(b>>8)&0xff]
	if a > 0 {
		a = 1
	}
	return ARGB1555(a<<15 | r<<10 | g<<5 | b)
}

// ARGB2222 is a 8-bit RGB color with alpha channel.
type ARGB2222 uint8

func (c ARGB2222) RGBA() (r, g, b, a uint32) {
	a = lut2to8[uint32(c&0b11_00_00_00)>>6]
	r = lut2to8[uint32(c&0b00_11_00_00)>>4]
	g = lut2to8[uint32(c&0b00_00_11_00)>>2]
	b = lut2to8[uint32(c&0b00_00_00_11)>>0]
	return
}

func ToARGB2222(c color.Color) ARGB2222 {
	if c, ok := c.(ARGB2222); ok {
		return c
	}
	r, g, b, a := c.RGBA()
	a = lut8to2[(a>>8)&0xff]
	r = lut8to2[(r>>8)&0xff]
	g = lut8to2[(g>>8)&0xff]
	b = lut8to2[(b>>8)&0xff]
	return ARGB2222(a<<6 | r<<4 | g<<2 | b)
}

// ARGB3332 is a 11-bit RGB color with alpha channel.
type ARGB3332 uint16

func (c ARGB3332) RGBA() (r, g, b, a uint32) {
	a = lut3to8[uint32(c&0b111_000_000_00)>>8]
	r = lut3to8[uint32(c&0b000_111_000_00)>>5]
	g = lut3to8[uint32(c&0b000_000_111_00)>>2]
	b = lut2to8[uint32(c&0b000_000_000_11)>>0]
	return
}

func ToARGB3332(c color.Color) ARGB3332 {
	if c, ok := c.(ARGB3332); ok {
		return c
	}
	r, g, b, a := c.RGBA()
	a = lut8to3[(a>>8)&0xff]
	r = lut8to3[(r>>8)&0xff]
	g = lut8to3[(g>>8)&0xff]
	b = lut8to2[(b>>8)&0xff]
	return ARGB3332(a<<8 | r<<5 | g<<2 | b)
}

// ARGB4444 is a 16-bit RGB color with alpha channel.
type ARGB4444 uint16

func (c ARGB4444) RGBA() (r, g, b, a uint32) {
	a = lut4to8[(c&0b1111_0000_0000_0000)>>12]
	r = lut4to8[(c&0b0000_1111_0000_0000)>>8]
	g = lut4to8[(c&0b0000_0000_1111_0000)>>4]
	b = lut4to8[(c&0b0000_0000_0000_1111)>>0]
	return
}

func ToARGB4444(c color.Color) ARGB4444 {
	if c, ok := c.(ARGB4444); ok {
		return c
	}
	r, g, b, a := c.RGBA()
	a = lut8to4[(a>>8)&0xff]
	r = lut8to4[(r>>8)&0xff]
	g = lut8to4[(g>>8)&0xff]
	b = lut8to4[(b>>8)&0xff]
	return ARGB4444(a<<12 | r<<8 | g<<4 | b)
}

// ARGB6565 is a 22-bit RGB color with alpha channel.
type ARGB6565 uint32

func (c ARGB6565) RGBA() (r, g, b, a uint32) {
	a = lut6to8[(c&0b111111_00000_000000_00000)>>17]
	r = lut5to8[(c&0b000000_11111_000000_00000)>>11]
	g = lut6to8[(c&0b000000_00000_111111_00000)>>5]
	b = lut5to8[(c&0b000000_00000_000000_11111)>>0]
	return
}

func ToARGB6565(c color.Color) ARGB6565 {
	if c, ok := c.(ARGB6565); ok {
		return c
	}
	r, g, b, a := c.RGBA()
	a = lut8to6[(a>>8)&0xff]
	r = lut8to5[(r>>8)&0xff]
	g = lut8to6[(g>>8)&0xff]
	b = lut8to5[(b>>8)&0xff]
	return ARGB6565(a<<17 | r<<11 | g<<5 | b)
}

// ARGB8888 is a 32-bit RGB color with alpha channel.
type ARGB8888 uint32

func (c ARGB8888) RGBA() (r, g, b, a uint32) {
	a = uint32(c&0xf000) >> 12
	r = uint32(c&0x0f00) >> 8
	g = uint32(c&0x00f0) >> 4
	b = uint32(c&0x000f) >> 0
	return
}

func ToARGB8888(c color.Color) ARGB8888 {
	if c, ok := c.(ARGB8888); ok {
		return c
	}
	r, g, b, a := c.RGBA()
	r >>= 8
	g >>= 8
	b >>= 8
	a >>= 8
	return ARGB8888((a&0xff)<<12 | (r&0xff)<<8 | (g&0xff)<<4 | (b & 0xff))
}

// ARGB2101010 is a 32-bit RGB color with alpha channel.
type ARGB2101010 uint32

func (c ARGB2101010) RGBA() (r, g, b, a uint32) {
	a = lut2to8[uint32(c>>30)&0x0003]
	r = uint32(c>>20) & 0x03ff
	r = r<<6 | r
	g = uint32(c>>10) & 0x03ff
	g = g<<6 | g
	b = uint32(c>>0) & 0x03ff
	b = b<<6 | b
	return
}

func ToARGB2101010(c color.Color) ARGB2101010 {
	if c, ok := c.(ARGB2101010); ok {
		return c
	}
	r, g, b, a := c.RGBA()
	a = (a >> 14) & 0x0003
	r = (r >> 6) & 0x03ff
	g = (g >> 6) & 0x03ff
	b = (b >> 6) & 0x03ff
	return ARGB2101010(a<<30 | r<<20 | g<<10 | b)
}

// XRGB1555 is a 16-bit RGB color without alpha channel.
type XRGB1555 uint16

func (c XRGB1555) RGBA() (r, g, b, a uint32) {
	r = lut5to8[(c&0b0_11111_00000_00000)>>10]
	g = lut5to8[(c&0b0_00000_11111_00000)>>5]
	b = lut5to8[(c&0b0_00000_00000_11111)>>0]
	a = 0xffff
	return
}

func ToXRGB1555(c color.Color) color.Color {
	if c, ok := c.(XRGB1555); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	r = lut8to5[(r>>8)&0xff]
	g = lut8to5[(g>>8)&0xff]
	b = lut8to5[(b>>8)&0xff]
	return XRGB1555(r<<10 | g<<5 | b)
}
