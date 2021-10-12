package xmlvector

import (
	"testing"

	"github.com/koykov/vector"
)

func TestProlog(t *testing.T) {
	vec := NewVector()
	t.Run("prolog", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeNull)
	})
}
