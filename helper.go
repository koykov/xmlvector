package xmlvector

import (
	"io"

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

func (h Helper) Beautify(w io.Writer, node *vector.Node) error {
	return beautify(w, node, 0)
}
