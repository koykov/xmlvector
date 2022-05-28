package xmlvector

import "errors"

var (
	ErrUnclosedProlog = errors.New("unclosed prolog instruction")
	ErrBadAttr        = errors.New("bad attribute")
	ErrNoRoot         = errors.New("no root tag")
	ErrUnclosedTag    = errors.New("unclosed tag")
)
