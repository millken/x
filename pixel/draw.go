package pixel

import (
	"encoding/binary"
	"image"
	"image/color"

	"github.com/millken/x/pixel/pixelcolor"
)

func Fill(b Image, c color.Color) {
	switch b := b.(type) {
	case *Bitmap:
		if pixelcolor.ToBit(c) {
			memset(b.Pix, 0xff)
		} else {
			memset(b.Pix, 0x00)
		}
	case *RGB565:
		var v [2]byte
		binary.BigEndian.PutUint16(v[:], uint16(pixelcolor.ToRGB565(c)))
		memsetSlice(b.Pix, v[:])
	case *RGB888:
		v := pixelcolor.ToRGB888(c)
		memsetSlice(b.Pix, []byte{v.R, v.G, v.B})
	default:
		// Naive and slow set pixel implementation.
		r := b.Bounds()
		for y := r.Min.Y; y < r.Max.Y; y++ {
			for x := r.Min.X; x < r.Max.X; x++ {
				b.Set(x, y, c)
			}
		}
	}
}

// FillRectangle draws a filled rectangle between the two points.
func FillRectangle(b Image, r image.Rectangle, c color.Color) {
	fillRectangle(b, r.Min.X, r.Min.Y, r.Max.X-r.Min.X, r.Max.Y-r.Min.Y, c)
}

func fillRectangle(b Image, x, y, width, height int, c color.Color) {
	switch b := b.(type) {
	case *Bitmap:
		switch b.Format {
		case MHMSBFormat:
			var v byte
			if pixelcolor.ToBit(c) {
				v = 0x01
			}
			for xx := x; xx < x+width; xx++ {
				offset := 7 - xx&0x7
				for yy := y; yy < y+height; yy++ {
					index := yy*b.Stride + xx
					b.Pix[index] = (b.Pix[index] & ^(0x01 << offset)) | v<<offset
				}
			}
			return

		case MVLSBFormat:
			var v byte
			if pixelcolor.ToBit(c) {
				v = 0x01
			}
			for ; height > 0; height, y = height-1, y+1 {
				index := (y>>3)*b.Stride + x
				offset := y & 0x07
				for w := 0; w < width; w++ {
					b.Pix[index+w] = (b.Pix[index+w] & ^(0x01 << offset)) | v<<offset
				}
			}
			return

		default:
		}

	case *RGB565:
		var (
			r  = b.Bounds()
			xe = min(x+width, r.Max.X)
			ye = min(y+height, r.Max.Y)
			c  = pixelcolor.ToRGB565(c)
			v  [2]byte
		)
		binary.BigEndian.PutUint16(v[:], uint16(c))
		for ; y < ye; y++ {
			o := y*b.Stride + x*2
			memsetSlice(b.Pix[o:o+(xe-x)*2], v[:])
		}

	case *RGB888:
		var (
			r  = b.Bounds()
			xe = min(x+width, r.Max.X)
			ye = min(y+height, r.Max.Y)
			v  = pixelcolor.ToRGB888(c)
		)
		for ; y < ye; y++ {
			o := y*b.Stride + x*3
			memsetSlice(b.Pix[o:o+(xe-x)*3], []byte{v.R, v.G, v.B})
		}

	default:
		// Naive pixel by pixel fill.
		var (
			r  = b.Bounds()
			xe = min(x+width, r.Max.X)
			ye = min(y+height, r.Max.Y)
		)
		for ; y < ye; y++ {
			for ; x < xe; x++ {
				b.Set(x, y, c)
			}
		}
	}
}

// Rectangle draws a rectangle outline between the two points.
func Rectangle(b Image, r image.Rectangle, c color.Color) {
	rectangle(b, r.Min.X, r.Min.Y, r.Dx(), r.Dy(), c, false)
}

func rectangle(b Image, x, y, width, height int, c color.Color, fill bool) {
	r := b.Bounds()
	if width < 1 || height < 1 || (x+width) <= 0 || (y+height) <= 0 || y >= r.Max.Y || x >= r.Max.X {
		return
	}
	var (
		xe = min(r.Max.X-1, x+width-1)
		ye = min(r.Max.Y-1, y+height-1)
	)
	x = max(x, 0)
	y = max(y, 0)
	if fill || height == 1 || width == 1 {
		fillRectangle(b, x, y, xe-x+1, ye-y+1, c)
	} else {
		fillRectangle(b, x, y, xe-x+1, 1, c)
		fillRectangle(b, x, y, 1, ye-y+1, c)
		fillRectangle(b, x, ye, xe-x+1, 1, c)
		fillRectangle(b, xe, y, 1, ye-y+1, c)
	}
}

// HLine draws a horizontal line at p with the given width.
func HLine(b Image, p image.Point, width int, c color.Color) {
	rectangle(b, p.X, p.Y, width, 1, c, false)
}

// VLine draws a vertical line at p with the given height.
func VLine(b Image, p image.Point, height int, c color.Color) {
	rectangle(b, p.X, p.Y, 1, height, c, false)
}

// Line draws a line between two points.
func Line(b Image, p0, p1 image.Point, c color.Color) {
	var (
		dx     = abs(p1.X - p0.X)
		dy     = abs(p1.Y - p0.Y)
		x, y   = p0.X, p0.Y
		sx, sy = -1, -1
	)
	if p0.X <= p1.X {
		sx = 1
	}
	if p0.Y <= p1.Y {
		sy = 1
	}
	if dx > dy {
		e := float64(dx) / 2.0
		for x != p1.X {
			b.Set(x, y, c)
			e -= float64(dy)
			if e < 0 {
				y += sy
				e += float64(dx)
			}
			x += sx
		}
	} else {
		e := float64(dy) / 2.0
		for y != p1.Y {
			b.Set(x, y, c)
			e -= float64(dx)
			if e < 0 {
				x += sx
				e += float64(dy)
			}
			y += sy
		}
	}
	b.Set(x, y, c)
}

// Circle draws a circle at point p with radius r.
func Circle(b Image, p image.Point, r int, c color.Color) {
	if r < 0 {
		return
	}
	// Bresenham midpoint circle algorithm.
	var (
		x1, y1, err = -r, 0, 2 - 2*r
		x, y        = p.X, p.Y
	)
	for {
		b.Set(x-x1, y+y1, c)
		b.Set(x-y1, y-x1, c)
		b.Set(x+x1, y-y1, c)
		b.Set(x+y1, y+x1, c)
		r = err
		if r > x1 {
			x1++
			err += x1*2 + 1
		}
		if r <= y1 {
			y1++
			err += y1*2 + 1
		}
		if x1 >= 0 {
			break
		}
	}
}
