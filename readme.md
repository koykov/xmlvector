# XML vector

XML parser based on [Vector API](https://github.com/koykov/vector) with minimum memory consumption.

### Usage

```go
src := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<俄语 լեզու="ռուսերեն">данные</俄语>`)
vec := xmlvector.Acquire()
defer xmlvector.Release(vec)
_ = vec.Parse(src)
fmt.Println(vec.Dot("俄语@լեզու")) // ռուսերեն
```
