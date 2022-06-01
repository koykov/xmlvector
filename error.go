package xmlvector

import "errors"

var (
	ErrBadAttr     = errors.New("bad attribute")
	ErrNoRoot      = errors.New("no root tag")
	ErrUnclosedTag = errors.New("unclosed tag")
	ErrUnexpToken  = errors.New("unexpected token")
)
