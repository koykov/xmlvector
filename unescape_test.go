package xmlvector

import (
	"bytes"
	"testing"
)

type stageUnescape struct {
	origin, expect []byte
}

var (
	stagesUnescape = map[string]*stageUnescape{
		"lt": {
			origin: []byte("ten &lt; twenty"),
			expect: []byte("ten < twenty"),
		},
		"gt": {
			origin: []byte("999999 &gt; 111111"),
			expect: []byte("999999 > 111111"),
		},
		"amp": {
			origin: []byte("hammer &amp; bolter"),
			expect: []byte("hammer & bolter"),
		},
		"apos": {
			origin: []byte("I&apos;d like to go"),
			expect: []byte("I'd like to go"),
		},
		"quot": {
			origin: []byte("Here is some &quot;Text&quot;"),
			expect: []byte(`Here is some "Text"`),
		},
		"unicode": {
			origin: []byte("company_name &#169; / &#x2122; brand_name"),
			expect: []byte("company_name © / ™ brand_name"),
		},
		"mixed": {
			origin: []byte(`&lt;sometext&gt;
Here is some &quot;Text&quot; that I&apos;d like to be &quot;escaped&quot; for XML
&amp; here is some Swedish: Tack. Vars&#229;god.
&lt;/sometext&gt;`),
			expect: []byte(`<sometext>
Here is some "Text" that I'd like to be "escaped" for XML
& here is some Swedish: Tack. Varsågod.
</sometext>`),
		},
	}
)

func getStageUnescape(key string) *stageUnescape {
	if st, ok := stagesUnescape[key]; ok {
		return st
	}
	return nil
}

func testUnescape(tb testing.TB, buf []byte) []byte {
	key := getTBName(tb)
	st := getStageUnescape(key)
	if st == nil {
		tb.Fatal("stage not found")
	}
	buf = append(buf[:0], st.origin...)
	buf = Unescape(buf)
	if !bytes.Equal(buf, st.expect) {
		tb.Error("unescape failed")
	}
	return buf
}

func benchUnescape(b *testing.B) {
	var buf []byte
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf = testUnescape(b, buf)
	}
}

func TestUnescape(t *testing.T) {
	t.Run("lt", func(t *testing.T) { testUnescape(t, nil) })
	t.Run("gt", func(t *testing.T) { testUnescape(t, nil) })
	t.Run("amp", func(t *testing.T) { testUnescape(t, nil) })
	t.Run("apos", func(t *testing.T) { testUnescape(t, nil) })
	t.Run("quot", func(t *testing.T) { testUnescape(t, nil) })
	t.Run("unicode", func(t *testing.T) { testUnescape(t, nil) })
	t.Run("mixed", func(t *testing.T) { testUnescape(t, nil) })
}

func BenchmarkUnescape(b *testing.B) {
	b.Run("lt", func(b *testing.B) { benchUnescape(b) })
	b.Run("gt", func(b *testing.B) { benchUnescape(b) })
	b.Run("amp", func(b *testing.B) { benchUnescape(b) })
	b.Run("apos", func(b *testing.B) { benchUnescape(b) })
	b.Run("quot", func(b *testing.B) { benchUnescape(b) })
	b.Run("unicode", func(b *testing.B) { benchUnescape(b) })
	b.Run("mixed", func(b *testing.B) { benchUnescape(b) })
}
