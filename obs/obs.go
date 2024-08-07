package obs

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type OP uint

const (
	OP_MUL_DIV OP = 1 << iota // 乘除法
	OP_MOD                    // 取模
	OP_XOR                    // 异或
	OP_SHIFT                  // 位移，根据输入的值，当值>128时，进行右移，否则左移
)

func opN(a, b uint8, op OP) (r1 string, r2 []uint8) {
	var shiftFn = func(a uint8) uint8 {
		var b uint8
		if a > 128 {
			b = maxRightShift(a)
		} else {
			b = maxLeftShift(a)
		}
		return rnd(b)
	}
	var shbit = shiftFn(a)
	var mid uint8
	r2 = []uint8{a}
	if op&OP_SHIFT != 0 {
		r2 = append(r2, shbit)
		if a > 128 {
			mid = a >> shbit
			r1 = "%s>>%s"
		} else {
			mid = a << shbit
			r1 = "%s<<%s"
		}
	}
	if op&OP_XOR != 0 {
		xor := uint8(rand.Intn(256))
		r2 = append(r2, xor)
		mid = mid ^ xor
		r1 = fmt.Sprintf("%s ^ %%s", r1)
	}
	if op&OP_MOD != 0 && mid > 128 {
		mod := uint8(rand.Intn(128))
		r2 = append(r2, mod)
		mid = mid % mod
		r1 = fmt.Sprintf("(%s) %%%% %%s", r1)
	}
	if mid == b {
		return
	}
	if mid > b {
		t1 := mid - b
		r2 = append(r2, t1)
		r1 = fmt.Sprintf("%s-%%s", r1)
	} else {
		t1 := b - mid
		r2 = append(r2, t1)
		r1 = fmt.Sprintf("%s+%%s", r1)
	}
	return
}

func maxLeftShift(value uint8) uint8 {
	maxShift := uint8(0)
	for i := uint8(1); i < 8; i++ {
		if uint16(value)<<i <= 255 {
			maxShift = i
		} else {
			break
		}
	}
	return maxShift
}

func maxRightShift(value uint8) uint8 {
	maxShift := uint8(0)
	for i := uint8(1); i < 8; i++ {
		if value>>i > 0 {
			maxShift = i
		} else {
			break
		}
	}
	return maxShift
}

func formatAsciiTable(ascii [0xff]byte) string {
	result := "[0xff]byte{"
	for i, v := range ascii {
		if i > 0 {
			result += ", "
		}
		result += fmt.Sprintf("0x%02x", v)
	}
	result += "}"
	return result
}
func shulffedAsciiTable() [0xff]byte {
	var ascii [0xff]byte
	for i := 0x00; i < 0xff; i++ {
		ascii[i] = byte(i)
	}

	// 洗牌函数
	shuffle := func(arr *[0xff]byte) {
		for i := 0x7F; i > 0; i-- {
			j := rnd(uint8(i))
			arr[i], arr[j] = arr[j], arr[i]
		}
	}

	// 多次洗牌
	for i := 0; i < 5; i++ {
		shuffle(&ascii)
	}

	// 应用不同的位运算进行进一步混淆
	for i := 0x00; i < 0xff; i++ {
		ascii[i] = (ascii[i] ^ 0x55) + 0x33
		ascii[i] = (ascii[i] << 1) | (ascii[i] >> 7)
	}

	return ascii
}
func rnd(n uint8) uint8 {
	rand.NewSource(time.Now().UnixNano())
	return uint8(rand.Intn(int(n) + 1))
}

func GenerateByTable(old string) string {
	var a = shulffedAsciiTable()
	var b = make(map[uint8]byte, 0xff)
	for i := uint8(0); i < 0xff; i++ {
		b[a[i]] = i
	}
	var c []string
	for _, v := range old {
		vv := uint8(v)
		r1 := rnd(0xff)
		o1, o2 := opN(r1, vv, OP_SHIFT|OP_XOR|OP_MOD)
		s0 := make([]any, len(o2))
		for i, v := range o2 {
			s0[i] = fmt.Sprintf("a[%d]", b[v])
		}
		s1 := fmt.Sprintf(o1, s0...)
		c = append(c, s1)
	}
	s := fmt.Sprintf(`var a=%s
	   return string([]byte{%s})`, formatAsciiTable(a), strings.Join(c, ", "))
	return s
}

func GenerateCode(old string) string {
	var a = shulffedAsciiTable()
	var b = make(map[uint8]byte, 0xff)
	for i := uint8(0); i < 0xff; i++ {
		b[a[i]] = i
	}
	var c []string
	for _, v := range old {
		vv := uint8(v)
		r1 := rnd(0xff)
		o1, o2 := opN(r1, vv, OP_SHIFT|OP_XOR|OP_MOD)
		s0 := make([]any, len(o2))
		for i, v := range o2 {
			s0[i] = fmt.Sprintf("0x%02x", a[b[v]]) //如果输出十六进制，可以使用%d
		}
		s1 := fmt.Sprintf(o1, s0...)
		c = append(c, s1)
	}
	s := fmt.Sprintf(`[]byte{%s}`, strings.Join(c, ", "))
	return s
}
