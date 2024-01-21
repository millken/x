package pixel

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func memset(b []byte, v byte) {
	l := len(b)
	if l == 0 {
		return
	}
	b[0] = v
	for i := 1; i < l; i <<= 1 {
		copy(b[i:], b[:i])
	}
}

func memsetSlice(b, v []byte) {
	l := len(b)
	i := copy(b, v)
	for ; i < l; i <<= 1 {
		copy(b[i:], b[:i])
	}
}

func zeroRange(b []byte, n, l int) {
	for i := n; i < l; i++ {
		b[i] = 0
	}
}
