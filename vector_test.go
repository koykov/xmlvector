package xmlvector

import (
	"testing"

	"github.com/koykov/vector"
)

func TestProlog(t *testing.T) {
	vec := NewVector()
	t.Run("prolog", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeObj)
		assertStr(t, vec, "@version", "1.1")
		assertStr(t, vec, "@encoding", "UTF-8")
	})
	t.Run("prologMiss", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeObj)
		assertStr(t, vec, "@version", "1.0")
	})
	t.Run("skipPI", func(t *testing.T) {
		assertParse(t, vec, nil, 0)
	})
	t.Run("skipDT", func(t *testing.T) {
		assertParse(t, vec, nil, 0)
	})
	t.Run("skipDTLocal", func(t *testing.T) {
		assertParse(t, vec, nil, 0)
	})
	t.Run("skipHeader", func(t *testing.T) {
		assertParse(t, vec, nil, 0)
	})
}
