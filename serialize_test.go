package xmlvector

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSerialize(t *testing.T) {
	vec := NewVector()
	t.Run("serialize/beautify", func(t *testing.T) {
		assertParse(t, vec, nil, 0)
		key := getTBName(t)
		st := getStage(key)
		require.NotNil(t, st, "stage not found")
		var buf bytes.Buffer
		_ = vec.Beautify(&buf)
		assert.Equal(t, st.fmt, buf.Bytes(), "beautify mismatch")
	})
	t.Run("serialize/marshal", func(t *testing.T) {
		assertParse(t, vec, nil, 0)
		key := getTBName(t)
		st := getStage(key)
		require.NotNil(t, st, "stage not found")
		var buf bytes.Buffer
		_ = vec.Marshal(&buf)
		assert.Equal(t, st.flat, buf.Bytes(), "marshal mismatch")
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
