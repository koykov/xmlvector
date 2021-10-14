package xmlvector

import "errors"

var (
	ErrUnclosedProlog = errors.New("unclosed prolog instruction")
	ErrBadAttr        = errors.New("bad attribute")
)
