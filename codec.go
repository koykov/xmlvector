package xmlvector

import (
	"io"

	"github.com/koykov/vector"
)

type Codec struct {
	vector.BaseCodec
}

func (h Codec) Decode(p *vector.Byteptr) ([]byte, error) {
	b := p.RawBytes()
	if p.CheckBit(flagEscape) {
		p.SetBit(flagEscape, false)
		b = Unescape(b)
		p.SetLen(len(b))
	}
	return b, nil
}

func (h Codec) Beautify(w io.Writer, node *vector.Node) error {
	return serialize(w, node, 0, true)
}

func (h Codec) Marshal(w io.Writer, node *vector.Node) error {
	return serialize(w, node, 0, false)
}
