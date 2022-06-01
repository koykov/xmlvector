package xmlvector

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/koykov/vector"
)

type stage struct {
	key string

	origin, fmt []byte
}

var (
	stages []stage
)

func init() {
	_ = filepath.Walk("testdata", func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".xml" && !strings.Contains(filepath.Base(path), ".fmt.xml") {
			st := stage{}
			st.key = strings.Replace(path, ".xml", "", 1)
			st.key = strings.Replace(st.key, "testdata/", "", 1)
			st.origin, _ = ioutil.ReadFile(path)
			if st.fmt, _ = ioutil.ReadFile(strings.Replace(path, ".xml", ".fmt.xml", 1)); len(st.fmt) > 0 {
				// st.fmt = bytealg.Trim(st.fmt, btNl)
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
	if st == nil {
		tb.Fatal("stage not found")
	}
	dst.Reset()
	err1 := dst.ParseCopy(st.origin)
	if err1 != nil {
		if err != nil {
			if err != err1 || dst.ErrorOffset() != errOffset {
				tb.Fatalf(`error mismatch, need "%s" at %d, got "%s" at %d`, err.Error(), errOffset, err1.Error(), dst.ErrorOffset())
			}
		} else {
			tb.Fatalf(`err "%s" caught by offset %d`, err1.Error(), dst.ErrorOffset())
		}
	}
	return dst
}

func assertType(tb testing.TB, vec *Vector, path string, typ vector.Type) {
	if typ1 := vec.Dot(path).Type(); typ1 != typ {
		tb.Error("type mismatch, need", typ, "got", typ1)
	}
}

func assertStr(tb testing.TB, vec *Vector, path, expect string, typ vector.Type) {
	var node *vector.Node
	if node = vec.Dot(path); node.Type() != typ {
		tb.Error("node type mismatch, need", typ, "got", node.Type())
		return
	}
	if v := node.String(); v != expect {
		tb.Error("node value mismatch, need", expect, "got", v)
	}
}
