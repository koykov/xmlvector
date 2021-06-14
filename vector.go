package xmlvector

import (
	"github.com/koykov/fastconv"
	"github.com/koykov/vector"
)

type Vector struct {
	vector.Vector
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

func NewVector() *Vector {
	vec := &Vector{}
	// todo implement helper.
	vec.Helper = nil
	return vec
}
