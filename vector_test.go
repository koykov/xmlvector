package xmlvector

import (
	"bytes"
	"testing"

	"github.com/koykov/vector"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProlog(t *testing.T) {
	vec := NewVector()
	t.Run("prolog/initial", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeObject)
		assertStr(t, vec, "@version", "1.1", vector.TypeAttribute)
		assertStr(t, vec, "@encoding", "UTF-8", vector.TypeAttribute)
		assertStr(t, vec, "version", "initial", vector.TypeString)
	})
	t.Run("prolog/missed", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeObject)
		assertStr(t, vec, "@version", "1.0", vector.TypeAttribute)
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
		assertStr(t, vec, "root", "Lorem ipsum dolor sit amet, consectetur adipiscing elit.", vector.TypeString)
	})
	t.Run("root/collapsed", func(t *testing.T) {
		assertParse(t, vec, nil, 0)
		assertType(t, vec, "root", vector.TypeObject)
	})
	t.Run("root/attr", func(t *testing.T) {
		assertParse(t, vec, nil, 0)
		assertStr(t, vec, "root@title", "Foo", vector.TypeAttribute)
		assertStr(t, vec, "root@descr", "Bar", vector.TypeAttribute)
		assertStr(t, vec, "root@arg0", "qwe", vector.TypeAttribute)
		assertStr(t, vec, "root@arg1", "15", vector.TypeAttribute)
	})
	t.Run("root/object", func(t *testing.T) {
		assertParse(t, vec, nil, 0)
		assertType(t, vec, "note", vector.TypeObject)
		assertStr(t, vec, "note.to", "Tove", vector.TypeString)
		assertStr(t, vec, "note.from", "Jani", vector.TypeString)
		assertStr(t, vec, "note.heading", "Reminder", vector.TypeString)
		assertStr(t, vec, "note.body", "Don't forget me this weekend!", vector.TypeString)
	})
	t.Run("root/array", func(t *testing.T) {
		assertParse(t, vec, nil, 0)
		assertType(t, vec, "CATALOG.CD", vector.TypeArray)
		vec.Dot("CATALOG.CD").Each(func(idx int, node *vector.Node) {
			switch idx {
			case 0:
				assert.Equal(t, "Empire Burlesque", node.Dot("TITLE").Value().String())
			case 1:
				assert.Equal(t, "Bonnie Tyler", node.Dot("ARTIST").Value().String())
			case 2:
				assert.Equal(t, "USA", node.Dot("COUNTRY").Value().String())
			case 3:
				assert.Equal(t, "Virgin records", node.Dot("COMPANY").Value().String())
			case 4:
				assert.Equal(t, "9.90", node.Dot("PRICE").Value().String())
			case 5:
				assert.Equal(t, "1998", node.Dot("YEAR").Value().String())
			}
		})
	})
	t.Run("root/mixed", func(t *testing.T) {
		assertParse(t, vec, nil, 0)
		assertType(t, vec, "result.listing", vector.TypeArray)
		vec.Dot("result.listing").Each(func(idx int, node *vector.Node) {
			switch idx {
			case 0:
				assert.Equal(t, "Poker US  ", node.Dot("@title").String())
			case 1:
				assert.Equal(t, "Pop Creative", node.Dot("@descr").String())
			case 2:
				assert.Equal(t, "p.npcta.xyz", node.Dot("@site").String())
			case 3:
				assert.Equal(t, "0.000018", node.Dot("@bid").String())
			case 4:
				assert.Equal(t, "https://g.co/tfXw4dB5w2M_4", node.Dot("@url").String())
				assert.Equal(t, "foobar", node.String())
			}
		})
	})
	t.Run("root/unicode", func(t *testing.T) {
		assertParse(t, vec, nil, 0)
		assertStr(t, vec, "俄语", "данные", vector.TypeObject)
		assertStr(t, vec, "俄语@լեզու", "ռուսերեն", vector.TypeAttribute)
	})
	t.Run("root/comment", func(t *testing.T) {
		assertParse(t, vec, nil, 0)
		assertStr(t, vec, "list.payload", "foobar", vector.TypeString)
	})
	t.Run("root/multi-comment", func(t *testing.T) {
		assertParse(t, vec, nil, 0)
		assertStr(t, vec, "list.title", "welcome", vector.TypeString)
		assertStr(t, vec, "list.payload", "foobar", vector.TypeString)
	})
	t.Run("root/cdata", func(t *testing.T) {
		assertParse(t, vec, nil, 0)
		assertStr(t, vec, "movie.raw", "Marquis Warren", vector.TypeString)
		assertStr(t, vec, "movie.cdata", `<strong>Main protagonist<strong> of "The Hateful Eight"`, vector.TypeString)
	})
	t.Run("root/sq-attr", func(t *testing.T) {
		assertParse(t, vec, nil, 0)
	})
	t.Run("root/fmt-comment", func(t *testing.T) {
		assertParse(t, vec, nil, 0)
	})
}

func TestReader(t *testing.T) {
	t.Run("reader", func(t *testing.T) {
		src := []byte(`<?xml version="1.1" encoding="UTF-8"?><root>Lorem ipsum dolor sit amet, consectetur adipiscing elit.</root>`)
		rdr := bytes.NewReader(src)
		vec := NewVector()
		err := vec.ParseReader(rdr)
		require.NoError(t, err)
	})
	t.Run("file", func(t *testing.T) {
		vec := NewVector()
		err := vec.ParseFile("testdata/root/object.xml")
		require.NoError(t, err)
	})
}

func BenchmarkProlog(b *testing.B) {
	b.Run("prolog/initial", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			assertType(b, vec, "", vector.TypeObject)
			assertStr(b, vec, "@version", "1.1", vector.TypeAttribute)
			assertStr(b, vec, "@encoding", "UTF-8", vector.TypeAttribute)
			assertStr(b, vec, "version", "initial", vector.TypeString)
		})
	})
	b.Run("prolog/missed", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			assertType(b, vec, "", vector.TypeObject)
			assertStr(b, vec, "@version", "1.0", vector.TypeAttribute)
		})
	})
	b.Run("prolog/skipPI", func(b *testing.B) { bench(b, func(vec *Vector) {}) })
	b.Run("prolog/skipDT", func(b *testing.B) { bench(b, func(vec *Vector) {}) })
	b.Run("prolog/skipDTLocal", func(b *testing.B) { bench(b, func(vec *Vector) {}) })
	b.Run("prolog/skipHeader", func(b *testing.B) { bench(b, func(vec *Vector) {}) })
}

func BenchmarkRoot(b *testing.B) {
	b.Run("root/static", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			assertStr(b, vec, "root", "Lorem ipsum dolor sit amet, consectetur adipiscing elit.", vector.TypeString)
		})
	})
	b.Run("root/collapsed", func(b *testing.B) {
		bench(b, func(vec *Vector) { assertType(b, vec, "root", vector.TypeObject) })
	})
	b.Run("root/attr", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			assertStr(b, vec, "root@title", "Foo", vector.TypeAttribute)
			assertStr(b, vec, "root@descr", "Bar", vector.TypeAttribute)
			assertStr(b, vec, "root@arg0", "qwe", vector.TypeAttribute)
			assertStr(b, vec, "root@arg1", "15", vector.TypeAttribute)
		})
	})
	b.Run("root/object", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			assertType(b, vec, "note", vector.TypeObject)
			assertStr(b, vec, "note.to", "Tove", vector.TypeString)
			assertStr(b, vec, "note.from", "Jani", vector.TypeString)
			assertStr(b, vec, "note.heading", "Reminder", vector.TypeString)
			assertStr(b, vec, "note.body", "Don't forget me this weekend!", vector.TypeString)
		})
	})
	b.Run("root/array", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			assertType(b, vec, "CATALOG.CD", vector.TypeArray)
			vec.Dot("CATALOG.CD").Each(func(idx int, node *vector.Node) {
				switch idx {
				case 0:
					assert.True(b, "Empire Burlesque" == node.Dot("TITLE").Value().String())
				case 1:
					assert.True(b, "Bonnie Tyler" == node.Dot("ARTIST").Value().String())
				case 2:
					assert.True(b, "USA" == node.Dot("COUNTRY").Value().String())
				case 3:
					assert.True(b, "Virgin records" == node.Dot("COMPANY").Value().String())
				case 4:
					assert.True(b, "9.90" == node.Dot("PRICE").Value().String())
				case 5:
					assert.True(b, "1998" == node.Dot("YEAR").Value().String())
				}
			})
		})
	})
	b.Run("root/mixed", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			assertType(b, vec, "result.listing", vector.TypeArray)
			vec.Dot("result.listing").Each(func(idx int, node *vector.Node) {
				switch idx {
				case 0:
					assert.True(b, "Poker US  " == node.Dot("@title").String())
				case 1:
					assert.True(b, "Pop Creative" == node.Dot("@descr").String())
				case 2:
					assert.True(b, "p.npcta.xyz" == node.Dot("@site").String())
				case 3:
					assert.True(b, "0.000018" == node.Dot("@bid").String())
				case 4:
					assert.True(b, "https://g.co/tfXw4dB5w2M_4" == node.Dot("@url").String())
					assert.True(b, "foobar" == node.String())
				}
			})
		})
	})
	b.Run("root/unicode", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			assertStr(b, vec, "俄语", "данные", vector.TypeObject)
			assertStr(b, vec, "俄语@լեզու", "ռուսերեն", vector.TypeAttribute)
		})
	})
	b.Run("root/comment", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			assertStr(b, vec, "list.payload", "foobar", vector.TypeString)
		})
	})
	b.Run("root/multi-comment", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			assertStr(b, vec, "list.title", "welcome", vector.TypeString)
			assertStr(b, vec, "list.payload", "foobar", vector.TypeString)
		})
	})
	b.Run("root/cdata", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			assertStr(b, vec, "movie.raw", "Marquis Warren", vector.TypeString)
			assertStr(b, vec, "movie.cdata", `<strong>Main protagonist<strong> of "The Hateful Eight"`, vector.TypeString)
		})
	})
}

func BenchmarkReader(b *testing.B) {
	b.Run("reader", func(b *testing.B) {
		b.ReportAllocs()
		src := []byte(`<?xml version="1.1" encoding="UTF-8"?><root>Lorem ipsum dolor sit amet, consectetur adipiscing elit.</root>`)
		var buf bytes.Buffer
		vec := NewVector()
		for i := 0; i < b.N; i++ {
			buf.Reset()
			vec.Reset()
			_, _ = buf.Write(src)
			_ = vec.ParseReader(&buf)
		}
	})
	b.Run("file", func(b *testing.B) {
		b.ReportAllocs()
		vec := NewVector()
		for i := 0; i < b.N; i++ {
			vec.Reset()
			_ = vec.ParseFile("testdata/root/object.xml")
		}
	})
}
