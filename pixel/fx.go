package pixel

/*
func Scroll(buffer Image, xstep, ystep int) {
	if xstep == 0 {
		if ystep > 0 {
			for n := 0; n < ystep; n++ {
				ScrollDown(buffer)
			}
			return
		}
		for n := 0; n < -ystep; n++ {
			ScrollUp(buffer)
		}
		return
	}

	var (
		sx, y, dx, dy, ex, ey int
		cx, cy                int
		b                     = buffer.Bounds()
	)
	if xstep < 0 {
		sx = 0
		ex = b.Dx() + xstep
		dx = 1
		cx = b.Dx() - 1
	} else {
		sx = b.Dx() - 1
		ex = xstep - 1
		dx = -1
	}
	if ystep < 0 {
		y = 0
		ey = b.Dy() + ystep
		dy = 1
		cy = b.Dy() - 1
	} else {
		y = b.Dy() - 1
		ey = ystep - 1
		dy = -1
	}
	log.Printf("scroll: clearx: %d, cleary: %d", cx, cy)
	for ; y != ey; y += dy {
		for x := sx; x != ex; x += dx {
			if x == cx || y == cy {
				buffer.Set(x, y, color.Black)
			} else {
				buffer.Set(x, y, buffer.At(x-xstep, y-ystep))
			}
		}
	}
}
*/

func ScrollUp(b Image) {
	switch b := b.(type) {
	case *Bitmap:
		copy(b.Pix, b.Pix[b.Stride:])
		l := len(b.Pix)
		zeroRange(b.Pix, l-b.Stride, l)
	case *RGB565:
		copy(b.Pix, b.Pix[b.Stride:])
		l := len(b.Pix)
		zeroRange(b.Pix, l-b.Stride, l)
	case *RGB888:
		copy(b.Pix, b.Pix[b.Stride:])
		l := len(b.Pix)
		zeroRange(b.Pix, l-b.Stride, l)
	}
}

func ScrollDown(b Image) {
	switch b := b.(type) {
	case *Bitmap:
		copy(b.Pix[b.Stride:], b.Pix)
		zeroRange(b.Pix, 0, b.Stride)
	case *RGB565:
		copy(b.Pix[b.Stride:], b.Pix)
		zeroRange(b.Pix, 0, b.Stride)
	case *RGB888:
		copy(b.Pix[b.Stride:], b.Pix)
		zeroRange(b.Pix, 0, b.Stride)
	}
}
