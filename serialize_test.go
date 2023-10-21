package xmlvector

import (
	"bytes"
	"testing"
)

func TestSerialize(t *testing.T) {
	vec := NewVector()
	t.Run("serialize/beautify", func(t *testing.T) {
		assertParse(t, vec, nil, 0)
		key := getTBName(t)
		st := getStage(key)
		var buf bytes.Buffer
		_ = vec.Beautify(&buf)
		if !bytes.Equal(buf.Bytes(), st.fmt) {
			t.FailNow()
		}
	})
	t.Run("serialize/marshal", func(t *testing.T) {
		assertParse(t, vec, nil, 0)
		key := getTBName(t)
		st := getStage(key)
		var buf bytes.Buffer
		_ = vec.Marshal(&buf)
		if !bytes.Equal(buf.Bytes(), st.flat) {
			t.FailNow()
		}
	})
}

func BenchmarkSerialize(b *testing.B) {
	b.Run("serialize/beautify", func(b *testing.B) {
		b.ReportAllocs()
		var buf bytes.Buffer
		vec := NewVector()
		for i := 0; i < b.N; i++ {
			assertParse(b, vec, nil, 0)
			_ = vec.Beautify(&buf)
			buf.Reset()
		}
	})
	b.Run("serialize/marshal", func(b *testing.B) {
		b.ReportAllocs()
		var buf bytes.Buffer
		vec := NewVector()
		for i := 0; i < b.N; i++ {
			assertParse(b, vec, nil, 0)
			_ = vec.Marshal(&buf)
			buf.Reset()
		}
	})
}
