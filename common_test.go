package xmlvector

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/koykov/bytealg"
	"github.com/koykov/vector"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stage struct {
	key string

	origin, fmt, flat []byte
}

var stages []stage

func init() {
	_ = filepath.Walk("testdata", func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".xml" && !strings.Contains(filepath.Base(path), ".fmt.xml") {
			st := stage{}
			st.key = strings.Replace(path, ".xml", "", 1)
			st.key = strings.Replace(st.key, "testdata/", "", 1)
			st.origin, _ = os.ReadFile(path)
			if st.fmt, _ = os.ReadFile(strings.Replace(path, ".xml", ".fmt.xml", 1)); len(st.fmt) > 0 {
				// st.fmt = bytealg.Trim(st.fmt, btNl)
			}
			if st.flat, _ = os.ReadFile(strings.Replace(path, ".xml", ".flat.xml", 1)); len(st.flat) > 0 {
				st.flat = bytealg.Trim(st.flat, btNl)
			}
			stages = append(stages, st)
		}
		return nil
	})
}

func getStage(key string) (st *stage) {
	for i := 0; i < len(stages); i++ {
		st1 := &stages[i]
		if st1.key == key {
			st = st1
		}
	}
	return st
}

func getTBName(tb testing.TB) string {
	key := tb.Name()
	return key[strings.Index(key, "/")+1:]
}

func bench(b *testing.B, fn func(vec *Vector)) {
	vec := NewVector()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vec = assertParse(b, vec, nil, 0)
		fn(vec)
	}
}

func assertParse(tb testing.TB, dst *Vector, err error, errOffset int) *Vector {
	key := getTBName(tb)
	st := getStage(key)
	require.NotNil(tb, st, "stage not found")
	dst.Reset()
	err1 := dst.ParseCopy(st.origin)
	if err1 != nil {
		if err != nil {
			assert.True(tb, errors.Is(err1, err), `error mismatch, need "%s" at %d, got "%s" at %d`, err.Error(), errOffset, err1.Error(), dst.ErrorOffset())
			assert.True(tb, errOffset == dst.ErrorOffset(), "error offset mismatch")
		} else {
			tb.Fatalf(`err "%s" caught by offset %d`, err1.Error(), dst.ErrorOffset())
		}
	}
	return dst
}

func assertType(tb testing.TB, vec *Vector, path string, typ vector.Type) {
	assert.True(tb, typ == vec.Dot(path).Type(), "type mismatch")
}

func assertStr(tb testing.TB, vec *Vector, path, expect string, typ vector.Type) {
	node := vec.Dot(path)
	if !assert.True(tb, typ == node.Type(), "node type mismatch") {
		return
	}
	assert.True(tb, expect == node.String(), "node value mismatch")
}
