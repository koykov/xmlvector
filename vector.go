package xmlvector

import (
	"github.com/koykov/fastconv"
	"github.com/koykov/vector"
)

const (
	flagEscape = 0
	flagAttr   = 1
	flagAlias  = 2
)

// Vector implements XML vector parser.
type Vector struct {
	vector.Vector
}

// NewVector makes new parser.
func NewVector() *Vector {
	vec := &Vector{}
	vec.Helper = helper
	return vec
}

// Parse parses source bytes.
func (vec *Vector) Parse(s []byte) error {
	return vec.parse(s, false)
}

// ParseStr parses source string.
func (vec *Vector) ParseStr(s string) error {
	return vec.parse(fastconv.S2B(s), false)
}

// ParseCopy copies source bytes and parse it.
func (vec *Vector) ParseCopy(s []byte) error {
	return vec.parse(s, true)
}

// ParseCopyStr copies source string and parse it.
func (vec *Vector) ParseCopyStr(s string) error {
	return vec.parse(fastconv.S2B(s), true)
}
