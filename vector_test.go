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
		assertStr(t, vec, "@version", "1.1", vector.TypeAttr)
		assertStr(t, vec, "@encoding", "UTF-8", vector.TypeAttr)
		assertStr(t, vec, "version", "initial", vector.TypeStr)
	})
	t.Run("prolog/missed", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeObj)
		assertStr(t, vec, "@version", "1.0", vector.TypeAttr)
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
		assertStr(t, vec, "root", "Lorem ipsum dolor sit amet, consectetur adipiscing elit.", vector.TypeStr)
	})
	t.Run("root/collapsed", func(t *testing.T) {
		assertParse(t, vec, nil, 0)
		assertType(t, vec, "root", vector.TypeObj)
	})
	t.Run("root/attr", func(t *testing.T) {
		assertParse(t, vec, nil, 0)
		assertStr(t, vec, "root@title", "Foo", vector.TypeAttr)
		assertStr(t, vec, "root@descr", "Bar", vector.TypeAttr)
		assertStr(t, vec, "root@arg0", "qwe", vector.TypeAttr)
		assertStr(t, vec, "root@arg1", "15", vector.TypeAttr)
	})
	t.Run("root/object", func(t *testing.T) {
		assertParse(t, vec, nil, 0)
		assertType(t, vec, "note", vector.TypeObj)
		assertStr(t, vec, "note.to", "Tove", vector.TypeStr)
		assertStr(t, vec, "note.from", "Jani", vector.TypeStr)
		assertStr(t, vec, "note.heading", "Reminder", vector.TypeStr)
		assertStr(t, vec, "note.body", "Don't forget me this weekend!", vector.TypeStr)
	})
	t.Run("root/array", func(t *testing.T) {
		assertParse(t, vec, nil, 0)
		assertType(t, vec, "CATALOG.CD", vector.TypeArr)
		vec.Dot("CATALOG.CD").Each(func(idx int, node *vector.Node) {
			switch idx {
			case 0:
				if node.Dot("TITLE").Value().String() != "Empire Burlesque" {
					t.FailNow()
				}
			case 1:
				if node.Dot("ARTIST").Value().String() != "Bonnie Tyler" {
					t.FailNow()
				}
			case 2:
				if node.Dot("COUNTRY").Value().String() != "USA" {
					t.FailNow()
				}
			case 3:
				if node.Dot("COMPANY").Value().String() != "Virgin records" {
					t.FailNow()
				}
			case 4:
				if node.Dot("PRICE").Value().String() != "9.90" {
					t.FailNow()
				}
			case 5:
				if node.Dot("YEAR").Value().String() != "1998" {
					t.FailNow()
				}
			}
		})
	})
	t.Run("root/mixed", func(t *testing.T) {
		assertParse(t, vec, nil, 0)
		assertType(t, vec, "result.listing", vector.TypeArr)
		vec.Dot("result.listing").Each(func(idx int, node *vector.Node) {
			switch idx {
			case 0:
				if node.Dot("@title").String() != "Poker US  " {
					t.FailNow()
				}
			case 1:
				if node.Dot("@descr").String() != "Pop Creative" {
					t.FailNow()
				}
			case 2:
				if node.Dot("@site").String() != "p.npcta.xyz" {
					t.FailNow()
				}
			case 3:
				if node.Dot("@bid").String() != "0.000018" {
					t.FailNow()
				}
			case 4:
				if node.Dot("@url").String() != "https://g.co/tfXw4dB5w2M_4" {
					t.FailNow()
				}
				if node.String() != "foobar" {
					t.FailNow()
				}
			}
		})
	})
}
