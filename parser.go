package xmlvector

import (
	"bytes"
	"fmt"

	"github.com/koykov/bytealg"
	"github.com/koykov/vector"
)

var (
	// Byte constants.
	bFmt = []byte(" \t\n\r")

	bPrologOpen  = []byte("<?xml")
	bPrologClose = []byte("?>")
)

// Main internal parser helper.
func (vec *Vector) parse(s []byte, copy bool) (err error) {
	s = bytealg.Trim(s, bFmt)
	if err = vec.SetSrc(s, copy); err != nil {
		return
	}

	offset := 0
	// Create root node and register it.
	root, i := vec.GetNode(0)

	// Parse source data.
	offset, err = vec.parseGeneric(0, offset, root)
	if err != nil {
		vec.SetErrOffset(offset)
		return err
	}
	vec.PutNode(i, root)

	// Check unparsed tail.
	if offset < vec.SrcLen() {
		vec.SetErrOffset(offset)
		return vector.ErrUnparsedTail
	}

	return
}

// Generic parser helper.
func (vec *Vector) parseGeneric(depth, offset int, node *vector.Node) (int, error) {
	var err error
	node.SetOffset(vec.Index.Len(depth))
	src := vec.Src()[offset:]
	fmt.Println(string(src[:5]))
	switch {
	case len(src) > 4 && bytes.Equal(src[:5], bPrologOpen):
		if len(src) > 5 && src[5] != ' ' {
			// Ignore processing instructions like `<?xml-stylesheet ... ?>`
			return 0, nil
		}
		lim := bytealg.IndexAt(src, bPrologClose, 5)
		if lim == -1 {
			return offset, ErrUnclosedProlog
		}
		prolog := src[6:lim]
		if err = vec.parseAttr(prolog, node); err != nil {
			return offset, err
		}
		offset += lim + len(bPrologClose)
	}
	return offset, err
}

func (vec *Vector) parseAttr(p []byte, node *vector.Node) error {
	_, _ = p, node
	return nil
}
