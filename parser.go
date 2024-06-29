package xmlvector

import (
	"bytes"
	"errors"

	"github.com/koykov/bytealg"
	"github.com/koykov/vector"
)

const (
	offsetVersionKey = 0
	offsetVersionVal = 7
	lenVersionKey    = 7
	lenVersionVal    = 3

	lenDTOpen = 9
	lenPIOpen = 16
)

var (
	// Byte constants.
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

	bCDATAOpen  = []byte("<![CDATA[")
	bCDATAClose = []byte("]]>")

	bCommentOpen  = []byte("<!--")
	bCommentClose = []byte("-->")

	// Default key-value pairs.
	bPairs = []byte("version1.0")

	errBadInit = errors.New("bad vector initialization, use xmlvector.NewVector() or xmlvector.Acquire()")
)

// Main internal parser helper.
func (vec *Vector) parse(s []byte, copy bool) (err error) {
	if !vec.init {
		err = errBadInit
		return
	}

	s = bytealg.TrimBytesFmt4(s)
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
		cn  *vector.Node
		cni int
	)
	node.SetOffset(vec.Index.Len(depth))
	if offset, err = vec.parseProlog(depth+1, offset, node); err != nil {
		return offset, err
	}
	if offset, eof = vec.skipHdr(offset); eof {
		return offset, vector.ErrUnexpEOF
	}
	if cn, cni, offset, err = vec.parseElement(depth+1, offset, node); err != nil {
		return offset, err
	}
	if cn != nil {
		vec.PutNode(cni, cn)
	}
	return offset, nil
}

// Parse prolog instruction `<?xml ... ?>`.
func (vec *Vector) parseProlog(depth, offset int, node *vector.Node) (int, error) {
	var (
		err error
		eof bool
	)
	node.SetOffset(vec.Index.Len(depth))
	src := vec.Src()[offset:]
	n := len(src)
	_ = src[n-1]
	if len(src) > 4 && bytes.Equal(src[:5], bPrologOpen) {
		offset = 5
		offset, _, err = vec.parseAttr(depth, offset, node)
	} else {
		attr, i := vec.GetChildWT(node, depth, vector.TypeAttr)
		attr.Key().Init(bPairs, offsetVersionKey, lenVersionKey)
		attr.Value().Init(bPairs, offsetVersionVal, lenVersionVal)
		vec.PutNode(i, attr)
		return offset, nil
	}
	if n-offset >= 2 && bytes.Equal(src[:offset+2], bPrologClose) {
		offset += 2
	}
	if offset, eof = skipCommentAndFmt(src, n, offset); eof {
		err = vector.ErrUnexpEOF
	}
	return offset, err
}

// Skip header part (doctype and processing instructions)
// PI == processing instructions
// eg: <?xml-stylesheet type="text/css" href="my-style.css"?>
func (vec *Vector) skipHdr(offset int) (int, bool) {
	src := vec.Src()
	n := len(src)
	_ = src[n-1]
	var dt, pi, eof bool
loop:
	srcc := src[offset:]
	// DT
	if len(srcc) < lenDTOpen {
		return offset, false
	}
	if dt = bytes.Equal(srcc[:lenDTOpen], bDocType); dt {
		// Check local DT.
		p0, p1, p2 := bytealg.IndexAtBytes(srcc, bDTElem, lenDTOpen), bytealg.IndexAtBytes(srcc, bDTPCDATA, lenDTOpen), bytealg.IndexAtBytes(srcc, bDTClose, lenDTOpen)
		if p0 != -1 && p1 > p0 && p2 > p1 {
			offset += p2 + 2
		} else if p := bytealg.IndexByteAtBytes(srcc, '>', lenDTOpen); p != -1 {
			// Check DTD file.
			offset += p + 1
		}
	}
	// PI
	if len(srcc) < lenPIOpen {
		return offset, false
	}
	if pi = bytes.Equal(srcc[:lenPIOpen], bPIOpen); pi {
		posClose := bytealg.IndexAtBytes(srcc, bPIClose, lenPIOpen)
		if posClose == -1 {
			return offset, true
		}
		offset += posClose + 2
	}
	if offset, eof = skipCommentAndFmt(src, n, offset); eof {
		return offset, true
	}
	if dt || pi {
		goto loop
	}
	return offset, false
}

// Try parse XML element.
func (vec *Vector) parseElement(depth, offset int, root *vector.Node) (*vector.Node, int, int, error) {
	var (
		err error
		p   int
		tag []byte

		eof, clp bool
	)
	src := vec.Src()
	srcp := vec.SrcAddr()
	n := len(src)
	_ = src[n-1]
	if src[offset] != '<' {
		return nil, -1, offset, ErrNoRoot
	}
	offset++
	if offset, eof = skipCommentAndFmt(src, n, offset); eof && depth > 1 {
		return nil, -1, offset, vector.ErrUnexpEOF
	}
	if p = bytealg.IndexAnyAtBytes(src, bAfterTag, offset); p == -1 {
		return nil, -1, offset, ErrUnclosedTag
	}
	p = skipNameTable(src, n, offset, p)

	node, i := vec.GetChildWT(root, depth, vector.TypeObj)
	node.SetOffset(vec.Index.Len(depth + 1))
	node.Key().InitRaw(srcp, offset, p-offset)

	tag = src[offset:p]
	offset = p

	if offset, eof = skipCommentAndFmt(src, n, offset); eof && depth > 1 {
		return node, i, offset, vector.ErrUnexpEOF
	}
	if c := src[offset]; c != '/' && c != '>' {
		if offset, clp, err = vec.parseAttr(depth+1, offset, node); err != nil {
			return node, i, offset, err
		}

		if clp {
			return node, i, offset, nil
		}

		if offset, err = vec.parseContent(depth, offset, node); err != nil {
			return node, i, offset, err
		}
		if offset, eof = skipCommentAndFmt(src, n, offset); eof && depth > 1 {
			return node, i, offset, vector.ErrUnexpEOF
		}
		if offset, err = skipCTag(src, n, offset, tag); err != nil {
			return node, i, offset, err
		}
		if offset, eof = skipCommentAndFmt(src, n, offset); eof && depth > 1 {
			return node, i, offset, vector.ErrUnexpEOF
		}
		return node, i, offset, nil
	}
	if src[offset] == '/' {
		if offset < n-1 && src[offset+1] == '>' {
			offset += 2
			return node, i, offset, nil
		} else {
			return node, i, offset, ErrUnclosedTag
		}
	}
	if src[offset] == '>' {
		offset++
		if offset, err = vec.parseContent(depth, offset, node); err != nil {
			return node, i, offset, err
		}
		if offset, err = skipCTag(src, n, offset, tag); err != nil {
			return node, i, offset, err
		}
		if offset, eof = skipCommentAndFmt(src, n, offset); eof && depth > 1 {
			return node, i, offset, vector.ErrUnexpEOF
		}
		return node, i, offset, nil
	}
	return node, i, offset, ErrUnclosedTag
}

// Try parse XML element content.
func (vec *Vector) parseContent(depth, offset int, root *vector.Node) (int, error) {
	var (
		p     int
		eof   bool
		cdata bool
		err   error
	)
	src := vec.Src()
	srcp := vec.SrcAddr()
	n := len(src)
	_ = src[n-1]
	if offset, eof = skipCommentAndFmt(src, n, offset); eof {
		return offset, vector.ErrUnexpEOF
	}
	offset, cdata = skipCDATA(src, n, offset)

	if src[offset] == '<' && !cdata {
		sl := n
		var (
			pn, cn *vector.Node
			cni    int
			arr    bool
		)
		for {
			if offset, eof = skipCommentAndFmt(src, n, offset); eof {
				return offset, vector.ErrUnexpEOF
			}
			if cn, cni, offset, err = vec.parseElement(depth+1, offset, root); err != nil {
				return offset, err
			}
			if cn != nil {
				vec.PutNode(cni, cn)
			}
			if offset, eof = skipCommentAndFmt(src, n, offset); eof {
				return offset, vector.ErrUnexpEOF
			}
			if !arr {
				if pn == nil && cn != nil {
					pn = cn
				} else if cn.KeyString() == pn.KeyString() {
					arr = true
				}
			}
			if sl == offset {
				break
			}
			if offset+1 < sl {
				if bytes.Equal(src[offset:offset+2], bCTag) {
					break
				}
			}
		}
		if arr {
			root.SetType(vector.TypeArr)
			*root.Value() = *pn.Key() // Use value as an alias for arrays.
			root.Value().SetBit(flagAlias, true)
		}
	} else {
		var d int
		if cdata {
			if p = bytealg.IndexAtBytes(src, bCDATAClose, offset); p == -1 {
				return offset, vector.ErrUnexpEOF
			}
			d = 3
		} else {
			if p = bytealg.IndexByteAtBytes(src, '<', offset); p == -1 {
				return offset, ErrUnclosedTag
			}
		}
		raw := src[offset:p]
		root.Value().InitRaw(srcp, offset, p-offset)
		root.Value().SetBit(flagEscape, vec.checkEscape(raw))
		if !root.Key().CheckBit(flagAttr) {
			root.SetType(vector.TypeStr)
		}
		offset = p + d
	}
	return offset, nil
}

// Try parse XML element attributes.
func (vec *Vector) parseAttr(depth, offset int, node *vector.Node) (int, bool, error) {
	var (
		err      error
		eof, clp bool
	)

	src := vec.Src()
	srcp := vec.SrcAddr()
	n := len(src)
	_ = src[n-1]

	for {
		if offset, eof = skipCommentAndFmt(src, n, offset); eof {
			return offset, clp, vector.ErrUnexpEOF
		}
		posName := offset
		posName1 := bytealg.IndexByteAtBytes(src, '=', offset)
		if posName1 == -1 {
			err = ErrBadAttr
			break
		}
		offset = posName1
		if offset, eof = skipCommentAndFmt(src, n, offset); eof {
			return offset, clp, vector.ErrUnexpEOF
		}
		offset++
		var c byte
		if c = src[offset]; c != '"' && c != '\'' {
			err = ErrBadAttr
			break
		}
		offset++
		posVal := offset
		posVal1 := bytealg.IndexByteAtBytes(src, c, offset)
		if posVal1 == -1 {
			err = ErrBadAttr
			break
		}

		attr, i := vec.GetChildWT(node, depth, vector.TypeAttr)
		attr.Key().InitRaw(srcp, posName, posName1-posName)
		val := src[posVal:posVal1]
		attr.Value().InitRaw(srcp, posVal, posVal1-posVal)
		attr.Value().SetBit(flagEscape, vec.checkEscape(val))
		vec.PutNode(i, attr)
		node.Key().SetBit(flagAttr, true)

		offset = posVal1 + 1
		if offset, eof = skipCommentAndFmt(src, n, offset); eof {
			return offset, clp, vector.ErrUnexpEOF
		}

		var brk bool
		b := src[offset]
		switch b {
		case '?', '/':
			offset++
			if src[offset] != '>' {
				return offset, clp, ErrUnexpToken
			}
			offset++
			brk, clp = true, true
		case '>':
			offset++
			brk = true
		}
		if brk {
			break
		}
	}
	return offset, clp, err
}

// Check p for escaped entities and glyphs.
func (vec *Vector) checkEscape(p []byte) bool {
	if len(p) == 0 {
		return false
	}
	offset := 0
loop:
	posAmp, posSC := bytealg.IndexByteAtBytes(p, '&', offset), bytealg.IndexByteAtBytes(p, ';', offset)
	if posAmp == -1 || posSC == -1 {
		return false
	}
	if posSC-posAmp >= 2 && posAmp-posSC < 5 {
		return true
	}
	offset = posSC
	goto loop
}
