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
		assertStrWT(t, vec, "@version", "1.1", false, vector.TypeAttr)
		assertStrWT(t, vec, "@encoding", "UTF-8", false, vector.TypeAttr)
	})
	t.Run("prologMiss", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeObj)
		assertStrWT(t, vec, "@version", "1.0", false, vector.TypeAttr)
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

func TestRoot(t *testing.T) {
	vec := NewVector()
	t.Run("rootStatic", func(t *testing.T) {
		assertParse(t, vec, nil, 0)
		assertStrWT(t, vec, "root", "Lorem ipsum dolor sit amet, consectetur adipiscing elit.", true, vector.TypeObj)
	})
	t.Run("rootCollapsed", func(t *testing.T) {
		assertParse(t, vec, nil, 0)
		assertType(t, vec, "root", vector.TypeObj)
	})
}
