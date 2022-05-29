package xmlvector

import (
	"testing"

	"github.com/koykov/vector"
)

func TestProlog(t *testing.T) {
	vec := NewVector()
	t.Run("prolog/initial", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeObj)
		assertStrWT(t, vec, "@version", "1.1", false, vector.TypeAttr)
		assertStrWT(t, vec, "@encoding", "UTF-8", false, vector.TypeAttr)
		assertStrWT(t, vec, "version", "initial", true, vector.TypeObj)
	})
	t.Run("prolog/missed", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeObj)
		assertStrWT(t, vec, "@version", "1.0", false, vector.TypeAttr)
	})
	t.Run("prolog/skipPI", func(t *testing.T) {
		assertParse(t, vec, nil, 0)
	})
	t.Run("prolog/skipDT", func(t *testing.T) {
		assertParse(t, vec, nil, 0)
	})
	t.Run("prolog/skipDTLocal", func(t *testing.T) {
		assertParse(t, vec, nil, 0)
	})
	t.Run("prolog/skipHeader", func(t *testing.T) {
		assertParse(t, vec, nil, 0)
	})
}

func TestRoot(t *testing.T) {
	vec := NewVector()
	t.Run("root/static", func(t *testing.T) {
		assertParse(t, vec, nil, 0)
		assertStrWT(t, vec, "root", "Lorem ipsum dolor sit amet, consectetur adipiscing elit.", true, vector.TypeObj)
	})
	t.Run("root/collapsed", func(t *testing.T) {
		assertParse(t, vec, nil, 0)
		assertType(t, vec, "root", vector.TypeObj)
	})
	t.Run("root/attr", func(t *testing.T) {
		assertParse(t, vec, nil, 0)
		assertStrWT(t, vec, "root@title", "Foo", false, vector.TypeAttr)
		assertStrWT(t, vec, "root@descr", "Bar", false, vector.TypeAttr)
		assertStrWT(t, vec, "root@arg0", "qwe", false, vector.TypeAttr)
		assertStrWT(t, vec, "root@arg1", "15", false, vector.TypeAttr)
	})
	t.Run("root/object", func(t *testing.T) {
		assertParse(t, vec, nil, 0)
		assertType(t, vec, "root", vector.TypeObj)
	})
}
