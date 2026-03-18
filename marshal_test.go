package marshal_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

// This file demonstrates Go JSON marshaling behavior across a matrix of:
//   - Data types: string, int, bool, float64, slice, map, struct
//   - Tag:        no tag, `json:"f"`, `json:"f,omitempty"`
//   - Receiver:   value, pointer (nil), pointer (non-nil)
//   - Value:      zero/empty, non-zero/non-empty

// ---------------------------------------------------------------------------
// Structs for each combination
// ---------------------------------------------------------------------------

// --- string ---

type StringNoTag struct {
	F string
}
type StringTag struct {
	F string `json:"f"`
}
type StringOmit struct {
	F string `json:"f,omitempty"`
}
type StringPtrNoTag struct {
	F *string
}
type StringPtrTag struct {
	F *string `json:"f"`
}
type StringPtrOmit struct {
	F *string `json:"f,omitempty"`
}

// --- int ---

type IntNoTag struct {
	F int
}
type IntTag struct {
	F int `json:"f"`
}
type IntOmit struct {
	F int `json:"f,omitempty"`
}
type IntPtrNoTag struct {
	F *int
}
type IntPtrTag struct {
	F *int `json:"f"`
}
type IntPtrOmit struct {
	F *int `json:"f,omitempty"`
}

// --- bool ---

type BoolNoTag struct {
	F bool
}
type BoolTag struct {
	F bool `json:"f"`
}
type BoolOmit struct {
	F bool `json:"f,omitempty"`
}
type BoolPtrNoTag struct {
	F *bool
}
type BoolPtrTag struct {
	F *bool `json:"f"`
}
type BoolPtrOmit struct {
	F *bool `json:"f,omitempty"`
}

// --- float64 ---

type Float64NoTag struct {
	F float64
}
type Float64Tag struct {
	F float64 `json:"f"`
}
type Float64Omit struct {
	F float64 `json:"f,omitempty"`
}
type Float64PtrNoTag struct {
	F *float64
}
type Float64PtrTag struct {
	F *float64 `json:"f"`
}
type Float64PtrOmit struct {
	F *float64 `json:"f,omitempty"`
}

// --- []string (slice) ---

type SliceNoTag struct {
	F []string
}
type SliceTag struct {
	F []string `json:"f"`
}
type SliceOmit struct {
	F []string `json:"f,omitempty"`
}
type SlicePtrNoTag struct {
	F *[]string
}
type SlicePtrTag struct {
	F *[]string `json:"f"`
}
type SlicePtrOmit struct {
	F *[]string `json:"f,omitempty"`
}

// --- map[string]string ---

type MapNoTag struct {
	F map[string]string
}
type MapTag struct {
	F map[string]string `json:"f"`
}
type MapOmit struct {
	F map[string]string `json:"f,omitempty"`
}
type MapPtrNoTag struct {
	F *map[string]string
}
type MapPtrTag struct {
	F *map[string]string `json:"f"`
}
type MapPtrOmit struct {
	F *map[string]string `json:"f,omitempty"`
}

// --- struct ---

type Inner struct {
	X int `json:"x"`
}

type StructNoTag struct {
	F Inner
}
type StructTag struct {
	F Inner `json:"f"`
}
type StructOmit struct {
	F Inner `json:"f,omitempty"`
}
type StructPtrNoTag struct {
	F *Inner
}
type StructPtrTag struct {
	F *Inner `json:"f"`
}
type StructPtrOmit struct {
	F *Inner `json:"f,omitempty"`
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func marshal(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return "ERROR: " + err.Error()
	}
	return string(b)
}

func strPtr(s string) *string       { return &s }
func intPtr(i int) *int             { return &i }
func boolPtr(b bool) *bool          { return &b }
func float64Ptr(f float64) *float64 { return &f }
func slicePtr(s []string) *[]string { return &s }
func mapPtr(m map[string]string) *map[string]string { return &m }
func innerPtr(i Inner) *Inner       { return &i }

// ---------------------------------------------------------------------------
// Matrix test
// ---------------------------------------------------------------------------

type row struct {
	label     string
	ptr       bool   // is the field a pointer?
	tagName   string // "no tag" | "tag" | "omitempty"
	valueKind string // "zero" | "non-zero" | "nil" etc.
	input     any    // value passed to json.Marshal
	want      string // expected JSON output
}

func TestMarshalMatrix(t *testing.T) {
	zeroStr := ""
	nonZeroStr := "hello"
	zeroInt := 0
	nonZeroInt := 42
	zeroBool := false
	nonZeroBool := true
	zeroFloat := 0.0
	nonZeroFloat := 3.14
	zeroSlice := []string{}
	nonZeroSlice := []string{"a", "b"}
	nilSlice := []string(nil)
	zeroMap := map[string]string{}
	nonZeroMap := map[string]string{"k": "v"}
	zeroInner := Inner{}
	nonZeroInner := Inner{X: 7}

	rows := []row{
		// ── string ─────────────────────────────────────────────────────────
		// Value types: zero value is always included unless omitempty.
		{"string", false, "no tag",    "zero",     StringNoTag{},                    `{"F":""}`},
		{"string", false, "no tag",    "non-zero",  StringNoTag{F: nonZeroStr},       `{"F":"hello"}`},
		{"string", false, "tag",       "zero",     StringTag{},                      `{"f":""}`},
		{"string", false, "tag",       "non-zero",  StringTag{F: nonZeroStr},         `{"f":"hello"}`},
		{"string", false, "omitempty", "zero",     StringOmit{},                     `{}`},           // empty string omitted
		{"string", false, "omitempty", "non-zero",  StringOmit{F: nonZeroStr},        `{"f":"hello"}`},
		// Pointer: nil → null (included); omitempty omits nil but keeps non-nil zero
		{"string", true,  "no tag",    "nil",      StringPtrNoTag{},                 `{"F":null}`},
		{"string", true,  "no tag",    "zero",     StringPtrNoTag{F: strPtr(zeroStr)},  `{"F":""}`},
		{"string", true,  "no tag",    "non-zero",  StringPtrNoTag{F: strPtr(nonZeroStr)}, `{"F":"hello"}`},
		{"string", true,  "tag",       "nil",      StringPtrTag{},                   `{"f":null}`},
		{"string", true,  "tag",       "zero",     StringPtrTag{F: strPtr(zeroStr)},    `{"f":""}`},
		{"string", true,  "tag",       "non-zero",  StringPtrTag{F: strPtr(nonZeroStr)}, `{"f":"hello"}`},
		{"string", true,  "omitempty", "nil",      StringPtrOmit{},                  `{}`},           // nil ptr omitted
		{"string", true,  "omitempty", "zero",     StringPtrOmit{F: strPtr(zeroStr)},   `{"f":""}`},  // ptr non-nil → NOT omitted even if ""
		{"string", true,  "omitempty", "non-zero",  StringPtrOmit{F: strPtr(nonZeroStr)}, `{"f":"hello"}`},

		// ── int ────────────────────────────────────────────────────────────
		{"int", false, "no tag",    "zero",     IntNoTag{},                   `{"F":0}`},
		{"int", false, "no tag",    "non-zero",  IntNoTag{F: nonZeroInt},      `{"F":42}`},
		{"int", false, "tag",       "zero",     IntTag{},                     `{"f":0}`},
		{"int", false, "tag",       "non-zero",  IntTag{F: nonZeroInt},        `{"f":42}`},
		{"int", false, "omitempty", "zero",     IntOmit{},                    `{}`},          // 0 omitted
		{"int", false, "omitempty", "non-zero",  IntOmit{F: nonZeroInt},       `{"f":42}`},
		{"int", true,  "no tag",    "nil",      IntPtrNoTag{},                `{"F":null}`},
		{"int", true,  "no tag",    "zero",     IntPtrNoTag{F: intPtr(zeroInt)},  `{"F":0}`},
		{"int", true,  "no tag",    "non-zero",  IntPtrNoTag{F: intPtr(nonZeroInt)}, `{"F":42}`},
		{"int", true,  "tag",       "nil",      IntPtrTag{},                  `{"f":null}`},
		{"int", true,  "tag",       "zero",     IntPtrTag{F: intPtr(zeroInt)},    `{"f":0}`},
		{"int", true,  "tag",       "non-zero",  IntPtrTag{F: intPtr(nonZeroInt)}, `{"f":42}`},
		{"int", true,  "omitempty", "nil",      IntPtrOmit{},                 `{}`},          // nil ptr omitted
		{"int", true,  "omitempty", "zero",     IntPtrOmit{F: intPtr(zeroInt)},   `{"f":0}`}, // ptr non-nil → NOT omitted
		{"int", true,  "omitempty", "non-zero",  IntPtrOmit{F: intPtr(nonZeroInt)}, `{"f":42}`},

		// ── bool ───────────────────────────────────────────────────────────
		{"bool", false, "no tag",    "zero",     BoolNoTag{},                   `{"F":false}`},
		{"bool", false, "no tag",    "non-zero",  BoolNoTag{F: nonZeroBool},     `{"F":true}`},
		{"bool", false, "tag",       "zero",     BoolTag{},                     `{"f":false}`},
		{"bool", false, "tag",       "non-zero",  BoolTag{F: nonZeroBool},       `{"f":true}`},
		{"bool", false, "omitempty", "zero",     BoolOmit{},                    `{}`},          // false omitted
		{"bool", false, "omitempty", "non-zero",  BoolOmit{F: nonZeroBool},      `{"f":true}`},
		{"bool", true,  "no tag",    "nil",      BoolPtrNoTag{},                `{"F":null}`},
		{"bool", true,  "no tag",    "zero",     BoolPtrNoTag{F: boolPtr(zeroBool)},   `{"F":false}`},
		{"bool", true,  "no tag",    "non-zero",  BoolPtrNoTag{F: boolPtr(nonZeroBool)}, `{"F":true}`},
		{"bool", true,  "tag",       "nil",      BoolPtrTag{},                  `{"f":null}`},
		{"bool", true,  "tag",       "zero",     BoolPtrTag{F: boolPtr(zeroBool)},     `{"f":false}`},
		{"bool", true,  "tag",       "non-zero",  BoolPtrTag{F: boolPtr(nonZeroBool)}, `{"f":true}`},
		{"bool", true,  "omitempty", "nil",      BoolPtrOmit{},                 `{}`},           // nil ptr omitted
		{"bool", true,  "omitempty", "zero",     BoolPtrOmit{F: boolPtr(zeroBool)},    `{"f":false}`}, // ptr non-nil → NOT omitted
		{"bool", true,  "omitempty", "non-zero",  BoolPtrOmit{F: boolPtr(nonZeroBool)}, `{"f":true}`},

		// ── float64 ────────────────────────────────────────────────────────
		{"float64", false, "no tag",    "zero",     Float64NoTag{},                    `{"F":0}`},
		{"float64", false, "no tag",    "non-zero",  Float64NoTag{F: nonZeroFloat},     `{"F":3.14}`},
		{"float64", false, "tag",       "zero",     Float64Tag{},                      `{"f":0}`},
		{"float64", false, "tag",       "non-zero",  Float64Tag{F: nonZeroFloat},       `{"f":3.14}`},
		{"float64", false, "omitempty", "zero",     Float64Omit{},                     `{}`},           // 0.0 omitted
		{"float64", false, "omitempty", "non-zero",  Float64Omit{F: nonZeroFloat},      `{"f":3.14}`},
		{"float64", true,  "no tag",    "nil",      Float64PtrNoTag{},                 `{"F":null}`},
		{"float64", true,  "no tag",    "zero",     Float64PtrNoTag{F: float64Ptr(zeroFloat)},    `{"F":0}`},
		{"float64", true,  "no tag",    "non-zero",  Float64PtrNoTag{F: float64Ptr(nonZeroFloat)}, `{"F":3.14}`},
		{"float64", true,  "tag",       "nil",      Float64PtrTag{},                   `{"f":null}`},
		{"float64", true,  "tag",       "zero",     Float64PtrTag{F: float64Ptr(zeroFloat)},      `{"f":0}`},
		{"float64", true,  "tag",       "non-zero",  Float64PtrTag{F: float64Ptr(nonZeroFloat)},  `{"f":3.14}`},
		{"float64", true,  "omitempty", "nil",      Float64PtrOmit{},                  `{}`},           // nil ptr omitted
		{"float64", true,  "omitempty", "zero",     Float64PtrOmit{F: float64Ptr(zeroFloat)},     `{"f":0}`},    // ptr non-nil → NOT omitted
		{"float64", true,  "omitempty", "non-zero",  Float64PtrOmit{F: float64Ptr(nonZeroFloat)}, `{"f":3.14}`},

		// ── []string (slice) ───────────────────────────────────────────────
		// nil slice marshals as null; empty slice marshals as [].
		// omitempty treats both nil and empty slice as "empty" and omits them.
		{"[]string", false, "no tag",    "nil",           SliceNoTag{F: nilSlice},                     `{"F":null}`},
		{"[]string", false, "no tag",    "empty",         SliceNoTag{F: zeroSlice},                    `{"F":[]}`},
		{"[]string", false, "no tag",    "non-empty",     SliceNoTag{F: nonZeroSlice},                 `{"F":["a","b"]}`},
		{"[]string", false, "tag",       "nil",           SliceTag{F: nilSlice},                       `{"f":null}`},
		{"[]string", false, "tag",       "empty",         SliceTag{F: zeroSlice},                      `{"f":[]}`},
		{"[]string", false, "tag",       "non-empty",     SliceTag{F: nonZeroSlice},                   `{"f":["a","b"]}`},
		{"[]string", false, "omitempty", "nil",           SliceOmit{F: nilSlice},                      `{}`},             // nil omitted
		{"[]string", false, "omitempty", "empty",         SliceOmit{F: zeroSlice},                     `{}`},             // empty slice also omitted
		{"[]string", false, "omitempty", "non-empty",     SliceOmit{F: nonZeroSlice},                  `{"f":["a","b"]}`},
		{"[]string", true,  "no tag",    "nil ptr",       SlicePtrNoTag{},                             `{"F":null}`},
		{"[]string", true,  "no tag",    "ptr→nil",       SlicePtrNoTag{F: slicePtr(nilSlice)},        `{"F":null}`},
		{"[]string", true,  "no tag",    "ptr→empty",     SlicePtrNoTag{F: slicePtr(zeroSlice)},       `{"F":[]}`},
		{"[]string", true,  "no tag",    "ptr→non-empty", SlicePtrNoTag{F: slicePtr(nonZeroSlice)},    `{"F":["a","b"]}`},
		{"[]string", true,  "tag",       "nil ptr",       SlicePtrTag{},                               `{"f":null}`},
		{"[]string", true,  "tag",       "ptr→empty",     SlicePtrTag{F: slicePtr(zeroSlice)},         `{"f":[]}`},
		{"[]string", true,  "tag",       "ptr→non-empty", SlicePtrTag{F: slicePtr(nonZeroSlice)},      `{"f":["a","b"]}`},
		{"[]string", true,  "omitempty", "nil ptr",       SlicePtrOmit{},                              `{}`},             // nil ptr omitted
		{"[]string", true,  "omitempty", "ptr→nil",       SlicePtrOmit{F: slicePtr(nilSlice)},         `{"f":null}`},     // ptr non-nil → NOT omitted, inner nil → null
		{"[]string", true,  "omitempty", "ptr→empty",     SlicePtrOmit{F: slicePtr(zeroSlice)},        `{"f":[]}`},       // ptr non-nil → NOT omitted
		{"[]string", true,  "omitempty", "ptr→non-empty", SlicePtrOmit{F: slicePtr(nonZeroSlice)},     `{"f":["a","b"]}`},

		// ── map[string]string ──────────────────────────────────────────────
		// Mirrors slice: nil map → null, empty map → {}, omitempty omits nil and empty.
		{"map", false, "no tag",    "nil",           MapNoTag{},                       `{"F":null}`},
		{"map", false, "no tag",    "empty",         MapNoTag{F: zeroMap},             `{"F":{}}`},
		{"map", false, "no tag",    "non-empty",     MapNoTag{F: nonZeroMap},          `{"F":{"k":"v"}}`},
		{"map", false, "tag",       "nil",           MapTag{},                         `{"f":null}`},
		{"map", false, "tag",       "empty",         MapTag{F: zeroMap},               `{"f":{}}`},
		{"map", false, "tag",       "non-empty",     MapTag{F: nonZeroMap},            `{"f":{"k":"v"}}`},
		{"map", false, "omitempty", "nil",           MapOmit{},                        `{}`},             // nil omitted
		{"map", false, "omitempty", "empty",         MapOmit{F: zeroMap},              `{}`},             // empty map also omitted
		{"map", false, "omitempty", "non-empty",     MapOmit{F: nonZeroMap},           `{"f":{"k":"v"}}`},
		{"map", true,  "no tag",    "nil ptr",       MapPtrNoTag{},                    `{"F":null}`},
		{"map", true,  "no tag",    "ptr→nil",       MapPtrNoTag{F: mapPtr(nil)},      `{"F":null}`},
		{"map", true,  "no tag",    "ptr→empty",     MapPtrNoTag{F: mapPtr(zeroMap)},  `{"F":{}}`},
		{"map", true,  "no tag",    "ptr→non-empty", MapPtrNoTag{F: mapPtr(nonZeroMap)}, `{"F":{"k":"v"}}`},
		{"map", true,  "tag",       "nil ptr",       MapPtrTag{},                      `{"f":null}`},
		{"map", true,  "tag",       "ptr→empty",     MapPtrTag{F: mapPtr(zeroMap)},    `{"f":{}}`},
		{"map", true,  "tag",       "ptr→non-empty", MapPtrTag{F: mapPtr(nonZeroMap)}, `{"f":{"k":"v"}}`},
		{"map", true,  "omitempty", "nil ptr",       MapPtrOmit{},                     `{}`},             // nil ptr omitted
		{"map", true,  "omitempty", "ptr→nil",       MapPtrOmit{F: mapPtr(nil)},       `{"f":null}`},     // ptr non-nil → NOT omitted
		{"map", true,  "omitempty", "ptr→empty",     MapPtrOmit{F: mapPtr(zeroMap)},   `{"f":{}}`},       // ptr non-nil → NOT omitted
		{"map", true,  "omitempty", "ptr→non-empty", MapPtrOmit{F: mapPtr(nonZeroMap)}, `{"f":{"k":"v"}}`},

		// ── struct ─────────────────────────────────────────────────────────
		// IMPORTANT: omitempty never omits a value-type struct — structs have no
		// "empty" concept in encoding/json. Only a nil *struct is omitted.
		{"struct", false, "no tag",    "zero",     StructNoTag{},                          `{"F":{"x":0}}`},
		{"struct", false, "no tag",    "non-zero",  StructNoTag{F: nonZeroInner},           `{"F":{"x":7}}`},
		{"struct", false, "tag",       "zero",     StructTag{},                            `{"f":{"x":0}}`},
		{"struct", false, "tag",       "non-zero",  StructTag{F: nonZeroInner},             `{"f":{"x":7}}`},
		{"struct", false, "omitempty", "zero",     StructOmit{},                           `{"f":{"x":0}}`}, // NOT omitted — zero struct is not "empty"
		{"struct", false, "omitempty", "non-zero",  StructOmit{F: nonZeroInner},            `{"f":{"x":7}}`},
		{"struct", true,  "no tag",    "nil",      StructPtrNoTag{},                       `{"F":null}`},
		{"struct", true,  "no tag",    "zero",     StructPtrNoTag{F: innerPtr(zeroInner)},  `{"F":{"x":0}}`},
		{"struct", true,  "no tag",    "non-zero",  StructPtrNoTag{F: innerPtr(nonZeroInner)}, `{"F":{"x":7}}`},
		{"struct", true,  "tag",       "nil",      StructPtrTag{},                         `{"f":null}`},
		{"struct", true,  "tag",       "zero",     StructPtrTag{F: innerPtr(zeroInner)},    `{"f":{"x":0}}`},
		{"struct", true,  "tag",       "non-zero",  StructPtrTag{F: innerPtr(nonZeroInner)}, `{"f":{"x":7}}`},
		{"struct", true,  "omitempty", "nil",      StructPtrOmit{},                        `{}`},           // nil ptr omitted
		{"struct", true,  "omitempty", "zero",     StructPtrOmit{F: innerPtr(zeroInner)},   `{"f":{"x":0}}`}, // ptr non-nil → NOT omitted
		{"struct", true,  "omitempty", "non-zero",  StructPtrOmit{F: innerPtr(nonZeroInner)}, `{"f":{"x":7}}`},
	}

	// Verify each row and print the matrix table.
	header := fmt.Sprintf("%-12s %-8s %-12s %-16s %-24s %s",
		"type", "pointer", "tag", "value", "want", "got")
	sep := strings.Repeat("-", 100)
	t.Log("\n" + sep)
	t.Log(header)
	t.Log(sep)

	for _, r := range rows {
		got := marshal(r.input)
		pass := got == r.want

		ptrStr := "no"
		if r.ptr {
			ptrStr = "yes"
		}
		status := "OK"
		if !pass {
			status = "FAIL"
		}
		t.Logf("%-12s %-8s %-12s %-16s %-24s %-24s %s",
			r.label, ptrStr, r.tagName, r.valueKind, r.want, got, status)

		if !pass {
			t.Errorf("%s | ptr=%v | %s | %s: want %s, got %s",
				r.label, r.ptr, r.tagName, r.valueKind, r.want, got)
		}
	}
	t.Log(sep)
}
