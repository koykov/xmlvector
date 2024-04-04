package xmlvector

import (
	"bytes"
	"unsafe"

	"github.com/koykov/bytealg"
)

// Skip close tag of XML element and return offset.
func skipCTag(src []byte, n, offset int, tag []byte) (int, error) {
	_ = src[n-1]
	if offset < n-2 && !bytes.Equal(src[offset:offset+2], bCTag) {
		return offset, ErrUnclosedTag
	}
	offset += 2
	offset += len(tag)
	if src[offset] != '>' {
		return offset, ErrUnclosedTag
	}
	return offset + 1, nil
}

// Skip comments.
func skipComment(src []byte, n, offset int) (int, bool) {
	_ = src[n-1]
	var p int
loop:
	if offset+4 >= n {
		return offset, false
	}
	if bytes.Equal(bCommentOpen, src[offset:offset+4]) {
		offset += 4
		if p = bytealg.IndexAtBytes(src, bCommentClose, offset); p == -1 {
			return offset, true
		}
		offset = p + 3
		goto loop
	}
	return offset, false
}

// Skip mixed formatting bytes and comments.
// See skipFmt and skipComment.
func skipCommentAndFmt(src []byte, n, offset int) (int, bool) {
	_ = src[n-1]
	var eof bool
	if eof = offset == n; eof {
		return offset, eof
	}
	poff := -1
	for poff != offset {
		poff = offset
		if offset, eof = skipFmtTable(src, n, offset); eof {
			return offset, true
		}
		if offset, eof = skipComment(src, n, offset); eof {
			return offset, true
		}
	}
	return offset, false
}

// Checks CDATA instruction.
func skipCDATA(src []byte, n, offset int) (int, bool) {
	_ = src[n-1]
	if offset+9 >= n {
		return offset, false
	}
	if bytes.Equal(bCDATAOpen, src[offset:offset+9]) {
		offset += 9
		return offset, true
	}
	return offset, false
}

// Skip element name before till first formatting byte.
// DEPRECATED: use skipNameTable instead.
func skipName(src []byte, n, offset, limit int) int {
	_ = src[n-1]
	for offset < limit {
		if c := src[offset]; c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			break
		}
		offset++
	}
	return offset
}

func skipNameTable(src []byte, n, offset, limit int) int {
	_ = src[n-1]
	for offset < limit {
		if skipTable[src[offset]] {
			break
		}
		offset++
	}
	return offset
}

// Skip formatting symbols like tabs, spaces, ...
//
// Returns index of next non-format symbol.
// DEPRECATED: use skipFmtTable instead.
func skipFmt(src []byte, n, offset int) (int, bool) {
	_ = src[n-1]
	if src[offset] > ' ' {
		return offset, false
	}
	_ = src[n-1]
	for ; offset < n; offset++ {
		c := src[offset]
		if c != ' ' && c != '\t' && c != '\n' && c != '\r' {
			return offset, false
		}
	}
	return offset, true
}

// Table based approach of skipFmt.
func skipFmtTable(src []byte, n, offset int) (int, bool) {
	_ = src[n-1]
	_ = skipTable[255]
	if n-offset > 512 {
		offset, _ = skipFmtBin8(src, n, offset)
	}
	for ; skipTable[src[offset]]; offset++ {
	}
	return offset, offset == n
}

// Binary based approach of skipFmt.
func skipFmtBin8(src []byte, n, offset int) (int, bool) {
	_ = src[n-1]
	_ = skipTable[255]
	if *(*uint64)(unsafe.Pointer(&src[offset])) == binNlSpace7 {
		offset += 8
		for offset < n && *(*uint64)(unsafe.Pointer(&src[offset])) == binSpace8 {
			offset += 8
		}
	}
	return offset, false
}

var (
	skipTable   = [256]bool{}
	binNlSpace7 uint64
	binSpace8   uint64
)

func init() {
	skipTable[' '] = true
	skipTable['\t'] = true
	skipTable['\n'] = true
	skipTable['\t'] = true

	binNlSpace7Bytes, binSpace8Bytes := []byte("\n       "), []byte("        ")
	binNlSpace7, binSpace8 = *(*uint64)(unsafe.Pointer(&binNlSpace7Bytes[0])), *(*uint64)(unsafe.Pointer(&binSpace8Bytes[0]))
}

var _ = skipFmt
