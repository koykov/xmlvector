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
		// return vector.ErrUnparsedTail
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
		if b := vec.SrcAt(offset); b == '?' || b == '>' {
			break
		}
	}
	return offset, err
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
