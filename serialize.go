package xmlvector

import (
	"io"

	"github.com/koykov/vector"
)

var (
	btSpace = []byte(` `)
	btEq    = []byte(`=`)
	btQuote = []byte(`"`)
	btTagO  = []byte(`<`)
	btTagC  = []byte(`>`)
	btNl    = []byte("\n")
	btTab   = []byte("\t")
)

func serialize(w io.Writer, node *vector.Node, depth int, indent bool) (err error) {
	_, _ = w.Write(bPrologOpen)
	_ = btAttr(w, node)
	_, _ = w.Write(bPrologClose)
	if indent {
		_, _ = w.Write(btNl)
	}

	node.Each(func(idx int, node *vector.Node) {
		if node.Type() != vector.TypeAttr {
			err = serialize1(w, node, depth+1, indent)
		}
	})
	return
}

func serialize1(w io.Writer, node *vector.Node, depth int, indent bool) (err error) {
	switch node.Type() {
	case vector.TypeObject, vector.TypeArray:
		if indent {
			writePad(w, depth-1)
		}
		_, _ = w.Write(btTagO)
		_, _ = w.Write(node.Key().Bytes())
		_ = btAttr(w, node)
		_, _ = w.Write(btTagC)

		if node.Value().Len() > 0 && !node.Value().CheckBit(flagAlias) {
			_, _ = w.Write(node.Value().Bytes())
		} else {
			if indent {
				_, _ = w.Write(btNl)
			}
			node.Each(func(idx int, node *vector.Node) {
				if node.Type() != vector.TypeAttribute {
					err = serialize1(w, node, depth+1, indent)
				}
			})
			if indent {
				writePad(w, depth-1)
			}
		}

		_, _ = w.Write(bCTag)
		_, _ = w.Write(node.Key().Bytes())
		_, _ = w.Write(btTagC)
		if indent {
			_, _ = w.Write(btNl)
		}
	default:
		_, _ = w.Write(node.Value().Bytes())
		if indent {
			_, _ = w.Write(btNl)
		}
	}
	return
}

func btAttr(w io.Writer, node *vector.Node) (err error) {
	node.Each(func(idx int, node *vector.Node) {
		if node.Type() == vector.TypeAttribute {
			_, _ = w.Write(btSpace)
			_, _ = w.Write(node.Key().Bytes())
			_, _ = w.Write(btEq)
			_, _ = w.Write(btQuote)
			_, _ = w.Write(node.Value().Bytes())
			_, _ = w.Write(btQuote)
		}
	})
	return
}

func writePad(w io.Writer, cnt int) {
	for i := 0; i < cnt; i++ {
		_, _ = w.Write(btTab)
	}
}
