package xmlvector

import (
	"github.com/koykov/vector"
)

type Helper struct{}

var (
	helper = Helper{}
)

func (h Helper) Indirect(p *vector.Byteptr) []byte {
	b := p.RawBytes()
	if p.CheckBit(flagEscape) {
		p.SetBit(flagEscape, false)
		b = Unescape(b)
		p.SetLen(len(b))
	}
	return b
}
