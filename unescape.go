package xmlvector

import (
	"bytes"
	"strconv"

	"github.com/koykov/bytealg"
	"github.com/koykov/fastconv"
)

var (
	beLt   = []byte("&lt;")
	beGt   = []byte("&gt;")
	beAmp  = []byte("&amp;")
	beApos = []byte("&apos;")
	beQuot = []byte("&quot;")
)

// Unescape byte array using itself as a destination.
func Unescape(p []byte) []byte {
	l, i, j, off := len(p), 0, 0, 0
	for {
		i = bytealg.IndexByteAtLUR(p, '&', off)
		if i < 0 || i+1 == l {
			break
		}
		off = i + 1
		j = bytealg.IndexByteAtLUR(p, ';', off)
		if j < 0 || j <= i {
			break
		}
		entity := p[i : j+1]
		if len(entity) < 4 {
			off = j + 1
			continue
		}
		switch {
		case bytes.Equal(entity, beLt):
			p[i] = '<'
			copy(p[i+1:], p[j+1:])
			l -= 3
		case bytes.Equal(entity, beGt):
			p[i] = '>'
			copy(p[i+1:], p[j+1:])
			l -= 3
		case bytes.Equal(entity, beAmp):
			p[i] = '&'
			copy(p[i+1:], p[j+1:])
			l -= 4
		case bytes.Equal(entity, beApos):
			p[i] = '\''
			copy(p[i+1:], p[j+1:])
			l -= 5
		case bytes.Equal(entity, beQuot):
			p[i] = '"'
			copy(p[i+1:], p[j+1:])
			l -= 5
		case entity[1] == '#':
			x := entity[2 : len(entity)-1]
			u, err := unescNum(x)
			if err != nil {
				i++
				continue
			}
			r := rune(u)
			s := string(r)
			z := len(s)
			copy(p[i:], s)
			copy(p[i+z:], p[j+1:])
			l -= len(entity) - z
		}

		p = p[:l]
	}
	return p
}

func unescNum(x []byte) (uint64, error) {
	if x[0] == 'x' {
		return strconv.ParseUint(fastconv.B2S(x[1:]), 16, 64)
	} else {
		return strconv.ParseUint(fastconv.B2S(x), 10, 64)
	}
}
