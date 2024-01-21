package pixeleffect

import (
	"image"
	"image/color"
	"image/draw"
	"time"

	"github.com/millken/x/pixel"
)

const (
	shiftRightRGBMaxSteps    = 8
	shiftRightRGB565MaxSteps = 6
	shiftRightRGB888MaxSteps = 8
)

func repeatWithin(duration time.Duration, n int, f func(int)) <-chan time.Time {
	signal := make(chan time.Time)
	if n == 0 {
		close(signal)
	} else {
		go func(signal chan<- time.Time, n int, f func(int)) {
			ticker := time.NewTicker(duration / time.Duration(n))
			for i := 0; i < n; i++ {
				f(i)
				signal <- <-ticker.C
			}
			ticker.Stop()
			close(signal)
		}(signal, n, f)
	}

	return signal
}

// FadeOutDither dims the colors by dithering, most useful for 1-bit bitmaps.
func FadeOutDither(duration time.Duration, im pixel.Image) <-chan time.Time {
	return repeatWithin(duration, len(ditherMasks), func(step int) {
		DitherStep(im, step)
	})
}

// FadeOutVertical draws horizontal raster lines to fade out the image.
func FadeOutHorizontal(duration time.Duration, im pixel.Image) <-chan time.Time {
	var (
		bounds = im.Bounds()
		width  = bounds.Dx()
		height = bounds.Dy()
		steps  = height / 2
	)
	if height&1 == 1 {
		steps++
	}
	return repeatWithin(duration, steps, func(y int) {
		var (
			offset0 = y * 2
			offset1 = height - offset0 - 1
		)
		pixel.Line(im, image.Point{Y: offset0}, image.Point{X: width, Y: offset0}, color.Black)
		pixel.Line(im, image.Point{Y: offset1}, image.Point{X: width, Y: offset1}, color.Black)
	})
}

// FadeOutVertical draws vertical raster lines to fade out the image.
func FadeOutVertical(duration time.Duration, im pixel.Image) <-chan time.Time {
	var (
		bounds = im.Bounds()
		width  = bounds.Dx()
		height = bounds.Dy()
		steps  = width / 2
	)
	if width&1 == 1 {
		steps++
	}
	return repeatWithin(duration, steps, func(x int) {
		var (
			offset0 = x * 2
			offset1 = width - offset0 - 1
		)
		pixel.Line(im, image.Point{X: offset0}, image.Point{X: offset0, Y: height}, color.Black)
		pixel.Line(im, image.Point{X: offset1}, image.Point{X: offset1, Y: height}, color.Black)
	})
}

func FadeOutDiagonal(duration time.Duration, im pixel.Image) <-chan time.Time {
	var (
		bounds     = im.Bounds()
		width      = bounds.Dx()
		height     = bounds.Dy()
		max, steps int
	)
	if width > height {
		max = width
	} else {
		max = height
	}
	steps = max / 2
	if max&1 == 1 {
		steps++
	}
	return repeatWithin(duration, steps, func(i int) {
		var (
			offset0 = i * 2
			offset1 = max - offset0 - 1
		)
		pixel.Line(im, image.Point{X: offset0}, image.Point{Y: offset0}, color.Black)
		pixel.Line(im, image.Point{X: offset1}, image.Point{Y: offset1}, color.Black)
	})
}

type ditherMask [4]byte

func (m ditherMask) At(x, y int) color.Color {
	if x < 0 {
		x = -x
	}
	if y < 0 {
		y = -y
	}
	if m[y&3]&(1<<(x&3)) == 0 {
		return color.Black
	}
	return color.Transparent
}

func (m ditherMask) Bounds() image.Rectangle {
	return image.Rectangle{
		Min: image.Point{X: -1e9, Y: -1e9},
		Max: image.Point{X: +1e9, Y: +1e9},
	}
}

func (m ditherMask) ColorModel() color.Model {
	return color.RGBAModel
}

var ditherMasks = [8]*ditherMask{
	{
		0b1111,
		0b1110,
		0b1111,
		0b1011,
	},
	{
		0b1111,
		0b1010,
		0b1111,
		0b1010,
	},
	{
		0b1101,
		0b1010,
		0b0111,
		0b1010,
	},
	{
		0b0101,
		0b1010,
		0b0101,
		0b1010,
	},
	{
		0b0101,
		0b1000,
		0b0101,
		0b0010,
	},
	{
		0b0101,
		0b0000,
		0b0101,
		0b0000,
	},
	{
		0b0100,
		0b0000,
		0b0001,
		0b0000,
	},
	{
		0b0000,
		0b0000,
		0b0000,
		0b0000,
	},
}

func DitherStep(im pixel.Image, step int) {
	if step < 0 || step >= len(ditherMasks) {
		// Nothing to do here.
		return
	}
	draw.Draw(im, im.Bounds(), ditherMasks[step], image.Point{}, draw.Over)
	//draw.DrawMask(im, im.Bounds(), im, image.Point{}, ditherMasks[step], image.Point{}, draw.Src)
}

// FadeOutShift dims the colors by bit shifting. 1-bit bitmaps are ignored.
func FadeOutShift(duration time.Duration, im pixel.Image) <-chan time.Time {
	switch im := im.(type) {
	case *pixel.Bitmap:
		// pointless
		signal := make(chan time.Time)
		close(signal)
		return signal

	case *pixel.RGB565:
		return repeatWithin(duration, shiftRightRGB565MaxSteps, func(_ int) { shiftRightRGB565(im) })

	case *pixel.RGB888:
		return repeatWithin(duration, shiftRightRGB888MaxSteps, func(_ int) { shiftRightRGB888(im) })

	default:
		return repeatWithin(duration, shiftRightRGBMaxSteps, func(_ int) { shiftRightRGB(im) })
	}
}

func ShiftRight(im pixel.Image) {
	switch im := im.(type) {
	case *pixel.Bitmap:
		// Pointless, ignored.

	case *pixel.RGB565:
		shiftRightRGB565(im)

	case *pixel.RGB888:
		shiftRightRGB888(im)

	default:
		shiftRightRGB(im)
	}
}

// Naive approach by converting the color to RGBA and shifting.
func shiftRightRGB(im pixel.Image) {
	b := im.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, b, a := im.At(x, y).RGBA()
			im.Set(x, y, color.RGBA{
				R: uint8(r>>8) >> 1,
				G: uint8(g>>8) >> 1,
				B: uint8(b>>8) >> 1,
				A: uint8(a>>8) >> 0,
			})
		}
	}
}

func shiftRightRGB565(im *pixel.RGB565) {
	for i, l := 0, len(im.Pix); i < l; i += 2 {
		var (
			hi      = im.Pix[i+0]
			lo      = im.Pix[i+1]
			r, g, b byte
		)
		r |= (hi & 0b11111000) >> 4       // RRRRR... -> ....RRRR
		g |= (hi & 0b00000111) << 2       // .....GGG -> ...GGG..
		g |= (lo & 0b11100000) >> 6       // GGG..... -> ......GG
		b |= (lo & 0b00011111) >> 1       // ...BBBBB -> ....BBBB
		im.Pix[i+0] = (r << 3) | (g >> 3) // .RRRR.GG
		im.Pix[i+1] = (g << 5) | b        // GGG.BBBB
	}
}

func shiftRightRGB888(im *pixel.RGB888) {
	for i, l := 0, len(im.Pix); i < l; i += 3 {
		im.Pix[i+0] >>= 1
		im.Pix[i+1] >>= 1
		im.Pix[i+2] >>= 1
	}
}
