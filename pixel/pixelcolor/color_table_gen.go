// +build ignore

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io"
	"math"
	"os"
)

func main() {
	outputFile := flag.String("output", "color_table.go", "output file")
	flag.Parse()

	b := new(bytes.Buffer)
	fmt.Fprintln(b, "package pixelcolor")
	fmt.Fprintln(b)
	fmt.Fprintln(b, "var (")

	for bits := 2; bits < 8; bits++ {
		generateLUT(b, bits)
	}

	/*
		// 2 to 8 bits
		fmt.Fprintln(b, "\t // lut2to8 contains pre-multiplied values")
		fmt.Fprint(b, "\tlut2to8 = [4]uint32{")
		for v := 0; v < 0b00000100; v++ {
			x := (v * 255 / 3)
			fmt.Fprintf(b, "%#04x,", x|x<<8)
		}
		fmt.Fprintln(b, "}")
		fmt.Fprintln(b, "\tlut8to2 = [256]uint32{")
		for v := 0; v < 0x100; v++ {
			fmt.Fprintf(b, "%#02x,", (v*3)/255)
			if v > 0 && (v%8) == 7 {
				fmt.Fprintln(b)
			}
		}
		fmt.Fprintln(b, "\t}")

		// 3 to 8 bits
		fmt.Fprintln(b, "\t // lut3to8 contains pre-multiplied values")
		fmt.Fprintln(b, "\tlut3to8 = [8]uint32{")
		for v := 0; v < 8; v++ {
			x := (v * 255 / 7)
			fmt.Fprintf(b, "%#04x,", x|x<<8)
		}
		fmt.Fprintln(b, "\n\t}")
		fmt.Fprintln(b, "\tlut8to3 = [256]uint32{")
		for v := 0; v < 0x100; v++ {
			fmt.Fprintf(b, "%#02x,", (v*7)/255)
			if v > 0 && (v%8) == 7 {
				fmt.Fprintln(b)
			}
		}
		fmt.Fprintln(b, "\t}")

		// 5 to 8 bits
		fmt.Fprintln(b, "\t // lut5to8 contains pre-multiplied values")
		fmt.Fprintln(b, "\tlut5to8 = [32]uint32{")
		for v := 0; v < 0b00100000; v++ {
			// x := (v * 255 / 31)
			// x := (v*527 + 23) >> 6
			x := int(math.Floor(float64(v)*255/31 + 0.5))
			fmt.Fprintf(b, "%#04x,", x|x<<8)
			if v > 0 && (v%8) == 7 {
				fmt.Fprintln(b)
			}
		}
		fmt.Fprintln(b, "\t}")
		fmt.Fprintln(b, "\tlut8to5 = [256]uint32{")
		for v := 0; v < 0x100; v++ {
			// x := (v*31)/255
			// x := (v*249 + 1014) >> 11
			x := int(math.Floor(float64(v)*31/255 + 0.5))
			fmt.Fprintf(b, "%#02x,", x)
			if v > 0 && (v%8) == 7 {
				fmt.Fprintln(b)
			}
		}
		fmt.Fprintln(b, "\t}")

		// 6 to 8 bits
		fmt.Fprintln(b, "\t // lut6to8 contains pre-multiplied values")
		fmt.Fprintln(b, "\tlut6to8 = [64]uint32{")
		for v := 0; v < 0b01000000; v++ {
			// x := (v * 255 / 63)
			x := (v*259 + 33) >> 6
			fmt.Fprintf(b, "%#04x,", x|x<<8)
			if v > 0 && (v%8) == 7 {
				fmt.Fprintln(b)
			}
		}
		fmt.Fprintln(b, "\t}")
		fmt.Fprintln(b, "\tlut8to6 = [256]uint32{")
		for v := 0; v < 0x100; v++ {
			// x := (v*63)/255
			x := (v*253 + 505) >> 10
			fmt.Fprintf(b, "%#02x,", x)
			if v > 0 && (v%8) == 7 {
				fmt.Fprintln(b)
			}
		}
		fmt.Fprintln(b, "\t}")
	*/

	fmt.Fprintln(b, ")")

	s, err := format.Source(b.Bytes())
	if err != nil {
		fatal(err)
	}

	f, err := os.Create(*outputFile)
	if err != nil {
		fatal("error opening", *outputFile+":", err)
	}
	defer f.Close()
	if _, err = f.Write(s); err != nil {
		fatal("error saving:", err)
	}
}

func fatal(v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
	os.Exit(1)
}

func generateLUT(w io.Writer, bits int) {
	max := 1 << bits
	fmt.Fprintf(w, "\t // lut%dto8 contains pre-multiplied values\n", bits)
	fmt.Fprintf(w, "\tlut%dto8 = [%d]uint32{\n", bits, max)
	for v := 0; v < 1<<bits; v++ {
		// x := (v * 255 / 31)
		// x := (v*527 + 23) >> 6
		x := int(math.Floor(float64(v)*255/float64(max-1) + 0.5))
		fmt.Fprintf(w, "%#04x,", x|x<<8)
		if v > 0 && (v%8) == 7 {
			fmt.Fprintln(w)
		}
	}
	fmt.Fprintln(w, "\t}")
	fmt.Fprintf(w, "\tlut8to%d = [256]uint32{\n", bits)
	for v := 0; v < 0x100; v++ {
		// x := (v*31)/255
		// x := (v*249 + 1014) >> 11
		x := int(math.Floor(float64(v)*float64(max-1)/255 + 0.5))
		fmt.Fprintf(w, "%#02x,", x)
		if v > 0 && (v%8) == 7 {
			fmt.Fprintln(w)
		}
	}
	fmt.Fprintln(w, "\t}")
}
