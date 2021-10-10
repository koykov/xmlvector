package xmlvector

import (
	"io"

	"github.com/koykov/fastconv"
	"github.com/koykov/vector"
)

type Vector struct {
	vector.Vector
}

func NewVector() *Vector {
	vec := &Vector{}
	// todo implement helper.
	vec.Helper = nil
	return vec
}

func (vec *Vector) Parse(s []byte) error {
	return vec.parse(s, false)
}

func (vec *Vector) ParseStr(s string) error {
	return vec.parse(fastconv.S2B(s), false)
}

func (vec *Vector) ParseCopy(s []byte) error {
	return vec.parse(s, true)
}

func (vec *Vector) ParseCopyStr(s string) error {
	return vec.parse(fastconv.S2B(s), true)
}

func (vec *Vector) Beautify(w io.Writer) error {
	r := vec.Root()
	return vec.beautify(w, r, 0)
}
