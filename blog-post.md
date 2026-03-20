# Go JSON Marshaling Behavior: The Complete Reference

`encoding/json` has a lot of edge cases. Field tags, pointer indirection, nil vs empty, `omitempty` — it's easy to get surprised. This post is a structured reference covering every combination across the common types.

The matrix axes are:
- **Type**: `string`, `int`, `bool`, `float64`, `[]string`, `map[string]string`, `struct`
- **Tag**: no tag, `json:"f"`, `json:"f,omitempty"`
- **Receiver**: value type vs pointer (`*T`)
- **Value**: zero/nil/empty vs non-zero/non-empty

Each table includes output for both `encoding/json` (v1) and `encoding/json/v2` (`GOEXPERIMENT=jsonv2`, Go 1.25+). Rows where they differ are marked **†**.

---

## Rules to know before reading the tables

**Field naming.** Without a tag, the field name is used as-is (`F`). With `json:"f"`, it becomes `f`.

**Nil pointer.** A nil `*T` always marshals as `null`, regardless of omitempty — unless the field has `omitempty`, in which case it is omitted entirely.

**`omitempty` definition (v1).** A field is omitted when its value is the zero value for its type: `""`, `0`, `false`, `0.0`, nil pointer, nil/empty slice, nil/empty map. Structs are never considered empty — `omitempty` has no effect on value-type struct fields.

**Non-nil pointer.** A non-nil pointer is never omitted by `omitempty` — the pointer itself is non-nil, so it is not the zero value. The pointed-to value is marshaled normally.

---

## string

| Pointer | Tag | Value | v1 | v2 |
|---------|-----|-------|----|----|
| no | no tag | `""` | `{"F":""}` | `{"F":""}` |
| no | no tag | `"hello"` | `{"F":"hello"}` | `{"F":"hello"}` |
| no | `json:"f"` | `""` | `{"f":""}` | `{"f":""}` |
| no | `json:"f"` | `"hello"` | `{"f":"hello"}` | `{"f":"hello"}` |
| no | `omitempty` | `""` | `{}` | `{}` |
| no | `omitempty` | `"hello"` | `{"f":"hello"}` | `{"f":"hello"}` |
| yes | no tag | nil | `{"F":null}` | `{"F":null}` |
| yes | no tag | ptr→`""` | `{"F":""}` | `{"F":""}` |
| yes | no tag | ptr→`"hello"` | `{"F":"hello"}` | `{"F":"hello"}` |
| yes | `json:"f"` | nil | `{"f":null}` | `{"f":null}` |
| yes | `json:"f"` | ptr→`""` | `{"f":""}` | `{"f":""}` |
| yes | `json:"f"` | ptr→`"hello"` | `{"f":"hello"}` | `{"f":"hello"}` |
| yes | `omitempty` | nil | `{}` | `{}` |
| yes | `omitempty` | ptr→`""` | `{"f":""}` | **`{}`** † |
| yes | `omitempty` | ptr→`"hello"` | `{"f":"hello"}` | `{"f":"hello"}` |

† v2 dereferences the pointer and omits the empty string behind it.

---

## int

| Pointer | Tag | Value | v1 | v2 |
|---------|-----|-------|----|----|
| no | no tag | `0` | `{"F":0}` | `{"F":0}` |
| no | no tag | `42` | `{"F":42}` | `{"F":42}` |
| no | `json:"f"` | `0` | `{"f":0}` | `{"f":0}` |
| no | `json:"f"` | `42` | `{"f":42}` | `{"f":42}` |
| no | `omitempty` | `0` | `{}` | **`{"f":0}`** † |
| no | `omitempty` | `42` | `{"f":42}` | `{"f":42}` |
| yes | no tag | nil | `{"F":null}` | `{"F":null}` |
| yes | no tag | ptr→`0` | `{"F":0}` | `{"F":0}` |
| yes | no tag | ptr→`42` | `{"F":42}` | `{"F":42}` |
| yes | `json:"f"` | nil | `{"f":null}` | `{"f":null}` |
| yes | `json:"f"` | ptr→`0` | `{"f":0}` | `{"f":0}` |
| yes | `json:"f"` | ptr→`42` | `{"f":42}` | `{"f":42}` |
| yes | `omitempty` | nil | `{}` | `{}` |
| yes | `omitempty` | ptr→`0` | `{"f":0}` | `{"f":0}` |
| yes | `omitempty` | ptr→`42` | `{"f":42}` | `{"f":42}` |

† v2 does not consider numeric zero to be "empty". Numbers are never omitted by v2's `omitempty`.

---

## bool

| Pointer | Tag | Value | v1 | v2 |
|---------|-----|-------|----|----|
| no | no tag | `false` | `{"F":false}` | `{"F":false}` |
| no | no tag | `true` | `{"F":true}` | `{"F":true}` |
| no | `json:"f"` | `false` | `{"f":false}` | `{"f":false}` |
| no | `json:"f"` | `true` | `{"f":true}` | `{"f":true}` |
| no | `omitempty` | `false` | `{}` | **`{"f":false}`** † |
| no | `omitempty` | `true` | `{"f":true}` | `{"f":true}` |
| yes | no tag | nil | `{"F":null}` | `{"F":null}` |
| yes | no tag | ptr→`false` | `{"F":false}` | `{"F":false}` |
| yes | no tag | ptr→`true` | `{"F":true}` | `{"F":true}` |
| yes | `json:"f"` | nil | `{"f":null}` | `{"f":null}` |
| yes | `json:"f"` | ptr→`false` | `{"f":false}` | `{"f":false}` |
| yes | `json:"f"` | ptr→`true` | `{"f":true}` | `{"f":true}` |
| yes | `omitempty` | nil | `{}` | `{}` |
| yes | `omitempty` | ptr→`false` | `{"f":false}` | `{"f":false}` |
| yes | `omitempty` | ptr→`true` | `{"f":true}` | `{"f":true}` |

† v2 does not consider `false` to be "empty". Booleans are never omitted by v2's `omitempty`.

---

## float64

| Pointer | Tag | Value | v1 | v2 |
|---------|-----|-------|----|----|
| no | no tag | `0.0` | `{"F":0}` | `{"F":0}` |
| no | no tag | `3.14` | `{"F":3.14}` | `{"F":3.14}` |
| no | `json:"f"` | `0.0` | `{"f":0}` | `{"f":0}` |
| no | `json:"f"` | `3.14` | `{"f":3.14}` | `{"f":3.14}` |
| no | `omitempty` | `0.0` | `{}` | **`{"f":0}`** † |
| no | `omitempty` | `3.14` | `{"f":3.14}` | `{"f":3.14}` |
| yes | no tag | nil | `{"F":null}` | `{"F":null}` |
| yes | no tag | ptr→`0.0` | `{"F":0}` | `{"F":0}` |
| yes | no tag | ptr→`3.14` | `{"F":3.14}` | `{"F":3.14}` |
| yes | `json:"f"` | nil | `{"f":null}` | `{"f":null}` |
| yes | `json:"f"` | ptr→`0.0` | `{"f":0}` | `{"f":0}` |
| yes | `json:"f"` | ptr→`3.14` | `{"f":3.14}` | `{"f":3.14}` |
| yes | `omitempty` | nil | `{}` | `{}` |
| yes | `omitempty` | ptr→`0.0` | `{"f":0}` | `{"f":0}` |
| yes | `omitempty` | ptr→`3.14` | `{"f":3.14}` | `{"f":3.14}` |

† Same as int/bool: v2 never omits numeric zero.

---

## []string (slice)

Nil and empty slices behave differently. An empty slice (`[]string{}`) marshals as `[]`. A nil slice marshals as `null` in v1, `[]` in v2. Both are omitted by `omitempty`.

| Pointer | Tag | Value | v1 | v2 |
|---------|-----|-------|----|----|
| no | no tag | nil | `{"F":null}` | **`{"F":[]}`** † |
| no | no tag | `[]` | `{"F":[]}` | `{"F":[]}` |
| no | no tag | `["a","b"]` | `{"F":["a","b"]}` | `{"F":["a","b"]}` |
| no | `json:"f"` | nil | `{"f":null}` | **`{"f":[]}`** † |
| no | `json:"f"` | `[]` | `{"f":[]}` | `{"f":[]}` |
| no | `json:"f"` | `["a","b"]` | `{"f":["a","b"]}` | `{"f":["a","b"]}` |
| no | `omitempty` | nil | `{}` | `{}` |
| no | `omitempty` | `[]` | `{}` | `{}` |
| no | `omitempty` | `["a","b"]` | `{"f":["a","b"]}` | `{"f":["a","b"]}` |
| yes | no tag | nil ptr | `{"F":null}` | `{"F":null}` |
| yes | no tag | ptr→nil | `{"F":null}` | **`{"F":[]}`** † |
| yes | no tag | ptr→`[]` | `{"F":[]}` | `{"F":[]}` |
| yes | no tag | ptr→`["a","b"]` | `{"F":["a","b"]}` | `{"F":["a","b"]}` |
| yes | `json:"f"` | nil ptr | `{"f":null}` | `{"f":null}` |
| yes | `json:"f"` | ptr→`[]` | `{"f":[]}` | `{"f":[]}` |
| yes | `json:"f"` | ptr→`["a","b"]` | `{"f":["a","b"]}` | `{"f":["a","b"]}` |
| yes | `omitempty` | nil ptr | `{}` | `{}` |
| yes | `omitempty` | ptr→nil | `{"f":null}` | **`{}`** † |
| yes | `omitempty` | ptr→`[]` | `{"f":[]}` | **`{}`** † |
| yes | `omitempty` | ptr→`["a","b"]` | `{"f":["a","b"]}` | `{"f":["a","b"]}` |

† v2 unifies nil and empty slices (both → `[]`), and dereferences pointers for `omitempty` checks.

---

## map[string]string

Same nil vs empty distinction as slices.

| Pointer | Tag | Value | v1 | v2 |
|---------|-----|-------|----|----|
| no | no tag | nil | `{"F":null}` | **`{"F":{}}`** † |
| no | no tag | `{}` | `{"F":{}}` | `{"F":{}}` |
| no | no tag | `{"k":"v"}` | `{"F":{"k":"v"}}` | `{"F":{"k":"v"}}` |
| no | `json:"f"` | nil | `{"f":null}` | **`{"f":{}}`** † |
| no | `json:"f"` | `{}` | `{"f":{}}` | `{"f":{}}` |
| no | `json:"f"` | `{"k":"v"}` | `{"f":{"k":"v"}}` | `{"f":{"k":"v"}}` |
| no | `omitempty` | nil | `{}` | `{}` |
| no | `omitempty` | `{}` | `{}` | `{}` |
| no | `omitempty` | `{"k":"v"}` | `{"f":{"k":"v"}}` | `{"f":{"k":"v"}}` |
| yes | no tag | nil ptr | `{"F":null}` | `{"F":null}` |
| yes | no tag | ptr→nil | `{"F":null}` | **`{"F":{}}`** † |
| yes | no tag | ptr→`{}` | `{"F":{}}` | `{"F":{}}` |
| yes | no tag | ptr→`{"k":"v"}` | `{"F":{"k":"v"}}` | `{"F":{"k":"v"}}` |
| yes | `json:"f"` | nil ptr | `{"f":null}` | `{"f":null}` |
| yes | `json:"f"` | ptr→`{}` | `{"f":{}}` | `{"f":{}}` |
| yes | `json:"f"` | ptr→`{"k":"v"}` | `{"f":{"k":"v"}}` | `{"f":{"k":"v"}}` |
| yes | `omitempty` | nil ptr | `{}` | `{}` |
| yes | `omitempty` | ptr→nil | `{"f":null}` | **`{}`** † |
| yes | `omitempty` | ptr→`{}` | `{"f":{}}` | **`{}`** † |
| yes | `omitempty` | ptr→`{"k":"v"}` | `{"f":{"k":"v"}}` | `{"f":{"k":"v"}}` |

† v2 unifies nil and empty maps (both → `{}`), and dereferences pointers for `omitempty` checks.

---

## struct

Structs are never "empty". `omitempty` on a value-type struct field is a no-op in both versions. Only a nil pointer to a struct is omitted.

| Pointer | Tag | Value | v1 | v2 |
|---------|-----|-------|----|----|
| no | no tag | `{X:0}` | `{"F":{"x":0}}` | `{"F":{"x":0}}` |
| no | no tag | `{X:7}` | `{"F":{"x":7}}` | `{"F":{"x":7}}` |
| no | `json:"f"` | `{X:0}` | `{"f":{"x":0}}` | `{"f":{"x":0}}` |
| no | `json:"f"` | `{X:7}` | `{"f":{"x":7}}` | `{"f":{"x":7}}` |
| no | `omitempty` | `{X:0}` | `{"f":{"x":0}}` | `{"f":{"x":0}}` |
| no | `omitempty` | `{X:7}` | `{"f":{"x":7}}` | `{"f":{"x":7}}` |
| yes | no tag | nil | `{"F":null}` | `{"F":null}` |
| yes | no tag | ptr→`{X:0}` | `{"F":{"x":0}}` | `{"F":{"x":0}}` |
| yes | no tag | ptr→`{X:7}` | `{"F":{"x":7}}` | `{"F":{"x":7}}` |
| yes | `json:"f"` | nil | `{"f":null}` | `{"f":null}` |
| yes | `json:"f"` | ptr→`{X:0}` | `{"f":{"x":0}}` | `{"f":{"x":0}}` |
| yes | `json:"f"` | ptr→`{X:7}` | `{"f":{"x":7}}` | `{"f":{"x":7}}` |
| yes | `omitempty` | nil | `{}` | `{}` |
| yes | `omitempty` | ptr→`{X:0}` | `{"f":{"x":0}}` | `{"f":{"x":0}}` |
| yes | `omitempty` | ptr→`{X:7}` | `{"f":{"x":7}}` | `{"f":{"x":7}}` |

No differences between v1 and v2 for structs.

---

## Quick-reference: omitempty behavior by type

| Type | Zero value | v1 omits? | v2 omits? |
|------|-----------|-----------|-----------|
| `string` | `""` | yes | yes |
| `int` / `float64` | `0` / `0.0` | yes | **no** |
| `bool` | `false` | yes | **no** |
| `[]T` | nil | yes | yes |
| `[]T` | `[]` (empty) | yes | yes |
| `map[K]V` | nil | yes | yes |
| `map[K]V` | `{}` (empty) | yes | yes |
| `struct` | `{...}` (any) | no | no |
| `*T` | nil | yes | yes |
| `*T` | non-nil ptr to zero value | no | **yes** (v2 dereferences) |
| `*T` | non-nil ptr to non-zero value | no | no |
