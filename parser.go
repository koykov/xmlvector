package xmlvector

import (
	"bytes"

	"github.com/koykov/bytealg"
	"github.com/koykov/vector"
)

const (
	flagBufSrc = 8

	offsetVersionKey = 0
	offsetVersionVal = 7
	lenVersionKey    = 7
	lenVersionVal    = 3

	lenDTOpen = 9
	lenPIOpen = 16
)

var (
	// Byte constants.
	bFmt = []byte(" \t\n\r")

	bPrologOpen  = []byte("<?xml")
	bPrologClose = []byte("?>")

	bDocType  = []byte("<!DOCTYPE")
	bDTElem   = []byte("<!ELEMENT")
	bDTPCDATA = []byte("#PCDATA")
	bDTClose  = []byte("]>")

	bPIOpen  = []byte("<?xml-stylesheet")
	bPIClose = []byte("?>")

	bAfterTag = []byte(" />")
	bCTag     = []byte("</")

	// Default key-value pairs.
	bPairs = []byte("version1.0")
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
	if offset, eof = vec.skipHdr(offset); eof {
		return offset, vector.ErrUnexpEOF
	}
	if offset, err = vec.parseRoot(depth+1, offset, node); err != nil {
		return offset, err
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
		attr, i := vec.GetChildWT(node, depth, vector.TypeAttr)
		attr.Key().Init(bPairs, offsetVersionKey, lenVersionKey)
		attr.Value().Init(bPairs, offsetVersionVal, lenVersionVal)
		vec.PutNode(i, attr)
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

// Skip header part (doctype and processing instructions)
// PI == processing instructions
// eg: <?xml-stylesheet type="text/css" href="my-style.css"?>
func (vec *Vector) skipHdr(offset int) (int, bool) {
	var dt, pi, eof bool
loop:
	src := vec.Src()[offset:]
	// DT
	if len(src) < lenDTOpen {
		return offset, false
	}
	if dt = bytes.Equal(src[:lenDTOpen], bDocType); dt {
		// Check local DT.
		p0, p1, p2 := bytealg.IndexAt(src, bDTElem, lenDTOpen), bytealg.IndexAt(src, bDTPCDATA, lenDTOpen), bytealg.IndexAt(src, bDTClose, lenDTOpen)
		if p0 != -1 && p1 > p0 && p2 > p1 {
			offset += p2 + 2
		} else if p := bytealg.IndexByteAtLR(src, '>', lenDTOpen); p != -1 {
			// Check DTD file.
			offset += p + 1
		}
	}
	// PI
	if len(src) < lenPIOpen {
		return offset, false
	}
	if pi = bytes.Equal(src[:lenPIOpen], bPIOpen); pi {
		posClose := bytealg.IndexAt(src, bPIClose, lenPIOpen)
		if posClose == -1 {
			return offset, true
		}
		offset += posClose + 2
	}
	if offset, eof = vec.skipFmt(offset); eof {
		return offset, true
	}
	if dt || pi {
		goto loop
	}
	return offset, false
}

func (vec *Vector) parseRoot(depth, offset int, node *vector.Node) (int, error) {
	var (
		err error
		p   int
		tag []byte
		eof bool
	)
	if vec.SrcAt(offset) != '<' {
		return offset, ErrNoRoot
	}
	offset++
	if p = bytealg.IndexAnyAt(vec.Src(), bAfterTag, offset); p == -1 {
		return offset, ErrUnclosedTag
	}

	root, i := vec.GetChildWT(node, depth, vector.TypeObj)
	defer vec.PutNode(i, root)
	root.Key().Init(vec.Src(), offset, p-offset)

	tag = vec.Src()[offset:p]
	offset = p

	switch vec.SrcAt(offset) {
	case ' ':
		offset++
		if offset, err = vec.parseAttr(depth+1, offset, root); err != nil {
			return offset, err
		}
	case '/':
		if offset < vec.SrcLen()-1 && vec.SrcAt(offset+1) == '>' {
			offset += 2
			return offset, nil
		}
		return offset, ErrUnclosedTag
	case '>':
		offset++
		if offset, err = vec.parseContent(depth, offset, root); err != nil {
			return offset, err
		}
		if offset, err = vec.mustCTag(offset, tag); err != nil {
			return offset, err
		}
		if offset, eof = vec.skipFmt(offset); eof && depth > 1 {
			return offset, vector.ErrUnexpEOF
		}
		return offset, nil
	}

	return offset, err
}

func (vec *Vector) parseContent(depth, offset int, node *vector.Node) (int, error) {
	var (
		p   int
		eof bool
		err error
	)
	if offset, eof = vec.skipFmt(offset); eof {
		return offset, vector.ErrUnexpEOF
	}
	if vec.SrcAt(offset) == '<' {
		sl := vec.SrcLen()
		for {
			if offset, eof = vec.skipFmt(offset); eof {
				return offset, vector.ErrUnexpEOF
			}
			if offset, err = vec.parseRoot(depth+1, offset, node); err != nil {
				return offset, err
			}
			if sl == offset {
				break
			}
			if offset+1 < sl {
				if bytes.Equal(vec.Src()[offset:offset+2], bCTag) {
					break
				}
			}
		}
	} else {
		if p = bytealg.IndexByteAtLR(vec.Src(), '<', offset); p == -1 {
			return offset, ErrUnclosedTag
		}
		raw := vec.Src()[offset:p]
		node.Value().Init(vec.Src(), offset, p-offset)
		node.Value().SetBit(flagBufSrc, vec.checkEscape(raw))
		offset = p
	}
	return offset, nil
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

		attr, i := vec.GetChildWT(node, depth, vector.TypeAttr)
		attr.Key().Init(vec.Src(), posName, posName1-posName)
		val := vec.Src()[posVal:posVal1]
		attr.Value().Init(vec.Src(), posVal, posVal1-posVal)
		attr.Value().SetBit(flagBufSrc, vec.checkEscape(val))
		vec.PutNode(i, attr)

		offset = posVal1 + 1
		if offset, eof = vec.skipFmt(offset); eof {
			return offset, vector.ErrUnexpEOF
		}

		var brk bool
		b := vec.SrcAt(offset)
		switch b {
		case '?', '/':
			offset++
			if vec.SrcAt(offset) != '>' {
				return offset, ErrUnexpToken
			}
			offset++
			brk = true
		case '>':
			brk = true
		}
		if brk {
			break
		}
	}
	return offset, err
}

func (vec *Vector) mustCTag(offset int, tag []byte) (int, error) {
	if offset < vec.SrcLen()-2 && !bytes.Equal(vec.Src()[offset:offset+2], bCTag) {
		return offset, ErrUnclosedTag
	}
	offset += 2
	offset += len(tag)
	if vec.SrcAt(offset) != '>' {
		return offset, ErrUnclosedTag
	}
	return offset + 1, nil
}

// Check p for escaped entities and glyphs.
func (vec *Vector) checkEscape(p []byte) bool {
	if len(p) == 0 {
		return false
	}
	offset := 0
loop:
	posAmp, posSC := bytealg.IndexByteAtLR(p, '&', offset), bytealg.IndexByteAtLR(p, ';', offset)
	if posAmp == -1 || posSC == -1 {
		return false
	}
	if posSC-posAmp >= 2 && posAmp-posSC < 5 {
		return true
	}
	offset = posSC
	goto loop
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
