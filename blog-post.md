# Go JSON Marshaling: encoding/json V1 vs V2

Go 1.25 ships `encoding/json/v2` as an experiment (`GOEXPERIMENT=jsonv2`). The API is mostly compatible, but several edge-case behaviors changed. This post shows the full matrix.

## The Three Breaking Differences

### 1. `omitempty` no longer omits numeric/boolean zero values

V1 treats `0`, `false`, `0.0` as "empty". V2 does not — only `""`, nil/empty slice, and nil/empty map are considered empty.

### 2. `omitempty` through pointers: V2 dereferences

V1 never omits a non-nil pointer, regardless of what it points to. V2 dereferences and checks the pointed-to value.

### 3. nil slice/map without `omitempty`: V2 encodes as `[]`/`{}`

V1 marshals nil slices and maps as `null`. V2 unifies nil and empty, producing `[]` or `{}`.

---

## Full Results Matrix

Rows marked **DIFF** have different output between V1 and V2.

### string

| Pointer | Tag | Value | V1 output | V2 output |
|---------|-----|-------|-----------|-----------|
| no | no tag | zero | `{"F":""}` | `{"F":""}` |
| no | no tag | non-zero | `{"F":"hello"}` | `{"F":"hello"}` |
| no | tag | zero | `{"f":""}` | `{"f":""}` |
| no | tag | non-zero | `{"f":"hello"}` | `{"f":"hello"}` |
| no | omitempty | zero | `{}` | `{}` |
| no | omitempty | non-zero | `{"f":"hello"}` | `{"f":"hello"}` |
| yes | no tag | nil | `{"F":null}` | `{"F":null}` |
| yes | no tag | ptr→`""` | `{"F":""}` | `{"F":""}` |
| yes | no tag | ptr→`"hello"` | `{"F":"hello"}` | `{"F":"hello"}` |
| yes | tag | nil | `{"f":null}` | `{"f":null}` |
| yes | tag | ptr→`""` | `{"f":""}` | `{"f":""}` |
| yes | tag | ptr→`"hello"` | `{"f":"hello"}` | `{"f":"hello"}` |
| yes | omitempty | nil | `{}` | `{}` |
| yes | omitempty | ptr→`""` | `{"f":""}` | **`{}`** ← **DIFF** |
| yes | omitempty | ptr→`"hello"` | `{"f":"hello"}` | `{"f":"hello"}` |

### int

| Pointer | Tag | Value | V1 output | V2 output |
|---------|-----|-------|-----------|-----------|
| no | no tag | zero | `{"F":0}` | `{"F":0}` |
| no | no tag | non-zero | `{"F":42}` | `{"F":42}` |
| no | tag | zero | `{"f":0}` | `{"f":0}` |
| no | tag | non-zero | `{"f":42}` | `{"f":42}` |
| no | omitempty | zero | `{}` | **`{"f":0}`** ← **DIFF** |
| no | omitempty | non-zero | `{"f":42}` | `{"f":42}` |
| yes | no tag | nil | `{"F":null}` | `{"F":null}` |
| yes | no tag | ptr→`0` | `{"F":0}` | `{"F":0}` |
| yes | no tag | ptr→`42` | `{"F":42}` | `{"F":42}` |
| yes | tag | nil | `{"f":null}` | `{"f":null}` |
| yes | tag | ptr→`0` | `{"f":0}` | `{"f":0}` |
| yes | tag | ptr→`42` | `{"f":42}` | `{"f":42}` |
| yes | omitempty | nil | `{}` | `{}` |
| yes | omitempty | ptr→`0` | `{"f":0}` | `{"f":0}` |
| yes | omitempty | ptr→`42` | `{"f":42}` | `{"f":42}` |

### bool

| Pointer | Tag | Value | V1 output | V2 output |
|---------|-----|-------|-----------|-----------|
| no | no tag | zero | `{"F":false}` | `{"F":false}` |
| no | no tag | non-zero | `{"F":true}` | `{"F":true}` |
| no | tag | zero | `{"f":false}` | `{"f":false}` |
| no | tag | non-zero | `{"f":true}` | `{"f":true}` |
| no | omitempty | zero | `{}` | **`{"f":false}`** ← **DIFF** |
| no | omitempty | non-zero | `{"f":true}` | `{"f":true}` |
| yes | no tag | nil | `{"F":null}` | `{"F":null}` |
| yes | no tag | ptr→`false` | `{"F":false}` | `{"F":false}` |
| yes | no tag | ptr→`true` | `{"F":true}` | `{"F":true}` |
| yes | tag | nil | `{"f":null}` | `{"f":null}` |
| yes | tag | ptr→`false` | `{"f":false}` | `{"f":false}` |
| yes | tag | ptr→`true` | `{"f":true}` | `{"f":true}` |
| yes | omitempty | nil | `{}` | `{}` |
| yes | omitempty | ptr→`false` | `{"f":false}` | `{"f":false}` |
| yes | omitempty | ptr→`true` | `{"f":true}` | `{"f":true}` |

### float64

| Pointer | Tag | Value | V1 output | V2 output |
|---------|-----|-------|-----------|-----------|
| no | no tag | zero | `{"F":0}` | `{"F":0}` |
| no | no tag | non-zero | `{"F":3.14}` | `{"F":3.14}` |
| no | tag | zero | `{"f":0}` | `{"f":0}` |
| no | tag | non-zero | `{"f":3.14}` | `{"f":3.14}` |
| no | omitempty | zero | `{}` | **`{"f":0}`** ← **DIFF** |
| no | omitempty | non-zero | `{"f":3.14}` | `{"f":3.14}` |
| yes | no tag | nil | `{"F":null}` | `{"F":null}` |
| yes | no tag | ptr→`0` | `{"F":0}` | `{"F":0}` |
| yes | no tag | ptr→`3.14` | `{"F":3.14}` | `{"F":3.14}` |
| yes | tag | nil | `{"f":null}` | `{"f":null}` |
| yes | tag | ptr→`0` | `{"f":0}` | `{"f":0}` |
| yes | tag | ptr→`3.14` | `{"f":3.14}` | `{"f":3.14}` |
| yes | omitempty | nil | `{}` | `{}` |
| yes | omitempty | ptr→`0` | `{"f":0}` | `{"f":0}` |
| yes | omitempty | ptr→`3.14` | `{"f":3.14}` | `{"f":3.14}` |

### []string (slice)

| Pointer | Tag | Value | V1 output | V2 output |
|---------|-----|-------|-----------|-----------|
| no | no tag | nil | `{"F":null}` | **`{"F":[]}`** ← **DIFF** |
| no | no tag | empty | `{"F":[]}` | `{"F":[]}` |
| no | no tag | non-empty | `{"F":["a","b"]}` | `{"F":["a","b"]}` |
| no | tag | nil | `{"f":null}` | **`{"f":[]}`** ← **DIFF** |
| no | tag | empty | `{"f":[]}` | `{"f":[]}` |
| no | tag | non-empty | `{"f":["a","b"]}` | `{"f":["a","b"]}` |
| no | omitempty | nil | `{}` | `{}` |
| no | omitempty | empty | `{}` | `{}` |
| no | omitempty | non-empty | `{"f":["a","b"]}` | `{"f":["a","b"]}` |
| yes | no tag | nil ptr | `{"F":null}` | `{"F":null}` |
| yes | no tag | ptr→nil | `{"F":null}` | **`{"F":[]}`** ← **DIFF** |
| yes | no tag | ptr→empty | `{"F":[]}` | `{"F":[]}` |
| yes | no tag | ptr→non-empty | `{"F":["a","b"]}` | `{"F":["a","b"]}` |
| yes | tag | nil ptr | `{"f":null}` | `{"f":null}` |
| yes | tag | ptr→empty | `{"f":[]}` | `{"f":[]}` |
| yes | tag | ptr→non-empty | `{"f":["a","b"]}` | `{"f":["a","b"]}` |
| yes | omitempty | nil ptr | `{}` | `{}` |
| yes | omitempty | ptr→nil | `{"f":null}` | **`{}`** ← **DIFF** |
| yes | omitempty | ptr→empty | `{"f":[]}` | **`{}`** ← **DIFF** |
| yes | omitempty | ptr→non-empty | `{"f":["a","b"]}` | `{"f":["a","b"]}` |

### map[string]string

| Pointer | Tag | Value | V1 output | V2 output |
|---------|-----|-------|-----------|-----------|
| no | no tag | nil | `{"F":null}` | **`{"F":{}}`** ← **DIFF** |
| no | no tag | empty | `{"F":{}}` | `{"F":{}}` |
| no | no tag | non-empty | `{"F":{"k":"v"}}` | `{"F":{"k":"v"}}` |
| no | tag | nil | `{"f":null}` | **`{"f":{}}`** ← **DIFF** |
| no | tag | empty | `{"f":{}}` | `{"f":{}}` |
| no | tag | non-empty | `{"f":{"k":"v"}}` | `{"f":{"k":"v"}}` |
| no | omitempty | nil | `{}` | `{}` |
| no | omitempty | empty | `{}` | `{}` |
| no | omitempty | non-empty | `{"f":{"k":"v"}}` | `{"f":{"k":"v"}}` |
| yes | no tag | nil ptr | `{"F":null}` | `{"F":null}` |
| yes | no tag | ptr→nil | `{"F":null}` | **`{"F":{}}`** ← **DIFF** |
| yes | no tag | ptr→empty | `{"F":{}}` | `{"F":{}}` |
| yes | no tag | ptr→non-empty | `{"F":{"k":"v"}}` | `{"F":{"k":"v"}}` |
| yes | tag | nil ptr | `{"f":null}` | `{"f":null}` |
| yes | tag | ptr→empty | `{"f":{}}` | `{"f":{}}` |
| yes | tag | ptr→non-empty | `{"f":{"k":"v"}}` | `{"f":{"k":"v"}}` |
| yes | omitempty | nil ptr | `{}` | `{}` |
| yes | omitempty | ptr→nil | `{"f":null}` | **`{}`** ← **DIFF** |
| yes | omitempty | ptr→empty | `{"f":{}}` | **`{}`** ← **DIFF** |
| yes | omitempty | ptr→non-empty | `{"f":{"k":"v"}}` | `{"f":{"k":"v"}}` |

### struct

Structs have no "empty" concept in either version — `omitempty` never omits a value-type struct. Only a nil `*struct` is omitted.

| Pointer | Tag | Value | V1 output | V2 output |
|---------|-----|-------|-----------|-----------|
| no | no tag | zero | `{"F":{"x":0}}` | `{"F":{"x":0}}` |
| no | no tag | non-zero | `{"F":{"x":7}}` | `{"F":{"x":7}}` |
| no | tag | zero | `{"f":{"x":0}}` | `{"f":{"x":0}}` |
| no | tag | non-zero | `{"f":{"x":7}}` | `{"f":{"x":7}}` |
| no | omitempty | zero | `{"f":{"x":0}}` | `{"f":{"x":0}}` |
| no | omitempty | non-zero | `{"f":{"x":7}}` | `{"f":{"x":7}}` |
| yes | no tag | nil | `{"F":null}` | `{"F":null}` |
| yes | no tag | ptr→zero | `{"F":{"x":0}}` | `{"F":{"x":0}}` |
| yes | no tag | ptr→non-zero | `{"F":{"x":7}}` | `{"F":{"x":7}}` |
| yes | tag | nil | `{"f":null}` | `{"f":null}` |
| yes | tag | ptr→zero | `{"f":{"x":0}}` | `{"f":{"x":0}}` |
| yes | tag | ptr→non-zero | `{"f":{"x":7}}` | `{"f":{"x":7}}` |
| yes | omitempty | nil | `{}` | `{}` |
| yes | omitempty | ptr→zero | `{"f":{"x":0}}` | `{"f":{"x":0}}` |
| yes | omitempty | ptr→non-zero | `{"f":{"x":7}}` | `{"f":{"x":7}}` |

---

## Summary of Differences

| Scenario | V1 | V2 |
|----------|----|----|
| `int` with `omitempty`, value `0` | `{}` | `{"f":0}` |
| `bool` with `omitempty`, value `false` | `{}` | `{"f":false}` |
| `float64` with `omitempty`, value `0.0` | `{}` | `{"f":0}` |
| `*string` with `omitempty`, ptr→`""` | `{"f":""}` | `{}` |
| `[]string` (nil), no omitempty | `{"F":null}` | `{"F":[]}` |
| `*[]string` (ptr→nil), no omitempty | `{"F":null}` | `{"F":[]}` |
| `*[]string` with `omitempty`, ptr→nil | `{"f":null}` | `{}` |
| `*[]string` with `omitempty`, ptr→`[]` | `{"f":[]}` | `{}` |
| `map` (nil), no omitempty | `{"F":null}` | `{"F":{}}` |
| `*map` (ptr→nil), no omitempty | `{"F":null}` | `{"F":{}}` |
| `*map` with `omitempty`, ptr→nil | `{"f":null}` | `{}` |
| `*map` with `omitempty`, ptr→`{}` | `{"f":{}}` | `{}` |

## Migration Checklist

- **`omitempty` on `int`/`float64`/`bool` fields**: If you relied on zero values being omitted, add a pointer (`*int`, etc.) or switch to explicit `null` handling.
- **`omitempty` on pointer fields**: V2 now dereferences. A `*string` pointing to `""` will be omitted by V2 but not V1. Audit all `*T` fields with `omitempty`.
- **nil slice/map fields**: Code that checks for `null` in JSON output will break. V2 always emits `[]`/`{}` for nil containers.
