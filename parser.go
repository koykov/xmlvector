package xmlvector

import (
	"bytes"

	"github.com/koykov/bytealg"
	"github.com/koykov/vector"
)

const (
	flagBufSrc = 8
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
	root, i := vec.GetNodeWT(0, vector.TypeObj)

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
	var (
		err error
		eof bool
	)
	node.SetOffset(vec.Index.Len(depth))
	if offset, err = vec.parseProlog(depth+1, offset, node); err != nil {
		return offset, err
	}
	if offset, eof = vec.skipPI(offset); eof {
		return offset, vector.ErrUnexpEOF
	}
	return offset, nil
}

func (vec *Vector) parseProlog(depth, offset int, node *vector.Node) (int, error) {
	var (
		err error
		eof bool
	)
	node.SetOffset(vec.Index.Len(depth))
	src := vec.Src()[offset:]
	if len(src) > 4 && bytes.Equal(src[:5], bPrologOpen) {
		offset = 5
		offset, err = vec.parseAttr(depth, offset, node)
	} else {
		return offset, nil
	}
	if vec.SrcLen()-offset >= 2 && bytes.Equal(vec.Src()[offset:offset+2], bPrologClose) {
		offset += 2
	}
	if offset, eof = vec.skipFmt(offset); eof {
		err = vector.ErrUnexpEOF
	}
	return offset, err
}

// PI == processing instructions
// eg: <?xml-stylesheet type="text/css" href="my-style.css"?>
func (vec *Vector) skipPI(offset int) (int, bool) {
	// todo implement me
	return offset, false
}

func (vec *Vector) parseAttr(depth, offset int, node *vector.Node) (int, error) {
	var (
		err error
		eof bool
	)
	for {
		if offset, eof = vec.skipFmt(offset); eof {
			return offset, vector.ErrUnexpEOF
		}
		posName := offset
		posName1 := bytealg.IndexByteAtLR(vec.Src(), '=', offset)
		if posName1 == -1 {
			err = ErrBadAttr
			break
		}
		offset = posName1
		if offset, eof = vec.skipFmt(offset); eof {
			return offset, vector.ErrUnexpEOF
		}
		offset++
		if vec.SrcAt(offset) != '"' {
			err = ErrBadAttr
			break
		}
		offset++
		posVal := offset
		posVal1 := bytealg.IndexByteAtLR(vec.Src(), '"', offset)
		if posVal1 == -1 {
			err = ErrBadAttr
			break
		}

		attr, i := vec.GetChildWT(node, depth, vector.TypeStr)
		name := vec.Src()[posName:posName1]
		boff := vec.BufLen()
		blim := posName1 - posName + 1
		vec.BufAppendStr("@")
		vec.BufAppend(name)
		attr.Key().Init(vec.Buf(), boff, blim)
		val := vec.Src()[posVal:posVal1]
		attr.Value().Init(vec.Src(), posVal, posVal1-posVal)
		// todo check val for escaped entities and glyphs (https://en.wikipedia.org/wiki/XML#Characters_and_escaping)
		_ = val
		// attr.Value().SetBit(flagBufSrc, true)
		vec.PutNode(i, attr)

		offset = posVal1 + 1
		if offset, eof = vec.skipFmt(offset); eof {
			return offset, vector.ErrUnexpEOF
		}
		if b := vec.SrcAt(offset); b == '?' || b == '>' {
			break
		}
	}
	return offset, err
}

func (vec *Vector) skipFmt(offset int) (int, bool) {
loop:
	if offset >= vec.SrcLen() {
		return offset, true
	}
	c := vec.SrcAt(offset)
	if c != bFmt[0] && c != bFmt[1] && c != bFmt[2] && c != bFmt[3] {
		return offset, false
	}
	offset++
	goto loop
}
