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
	label    string
	ptr      bool   // is the field a pointer?
	omit     bool   // has omitempty?
	tagName  string // "no tag" | "tag" | "omitempty"
	valueKind string // "zero" | "non-zero" | "nil" (only for ptr)
	result   string
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
		{"string", false, false, "no tag", "zero", marshal(StringNoTag{})},
		{"string", false, false, "no tag", "non-zero", marshal(StringNoTag{F: nonZeroStr})},
		{"string", false, false, "tag", "zero", marshal(StringTag{})},
		{"string", false, false, "tag", "non-zero", marshal(StringTag{F: nonZeroStr})},
		{"string", false, true, "omitempty", "zero", marshal(StringOmit{})},
		{"string", false, true, "omitempty", "non-zero", marshal(StringOmit{F: nonZeroStr})},
		{"string", true, false, "no tag", "nil", marshal(StringPtrNoTag{})},
		{"string", true, false, "no tag", "zero", marshal(StringPtrNoTag{F: strPtr(zeroStr)})},
		{"string", true, false, "no tag", "non-zero", marshal(StringPtrNoTag{F: strPtr(nonZeroStr)})},
		{"string", true, false, "tag", "nil", marshal(StringPtrTag{})},
		{"string", true, false, "tag", "zero", marshal(StringPtrTag{F: strPtr(zeroStr)})},
		{"string", true, false, "tag", "non-zero", marshal(StringPtrTag{F: strPtr(nonZeroStr)})},
		{"string", true, true, "omitempty", "nil", marshal(StringPtrOmit{})},
		{"string", true, true, "omitempty", "zero", marshal(StringPtrOmit{F: strPtr(zeroStr)})},
		{"string", true, true, "omitempty", "non-zero", marshal(StringPtrOmit{F: strPtr(nonZeroStr)})},

		// ── int ────────────────────────────────────────────────────────────
		{"int", false, false, "no tag", "zero", marshal(IntNoTag{})},
		{"int", false, false, "no tag", "non-zero", marshal(IntNoTag{F: nonZeroInt})},
		{"int", false, false, "tag", "zero", marshal(IntTag{})},
		{"int", false, false, "tag", "non-zero", marshal(IntTag{F: nonZeroInt})},
		{"int", false, true, "omitempty", "zero", marshal(IntOmit{})},
		{"int", false, true, "omitempty", "non-zero", marshal(IntOmit{F: nonZeroInt})},
		{"int", true, false, "no tag", "nil", marshal(IntPtrNoTag{})},
		{"int", true, false, "no tag", "zero", marshal(IntPtrNoTag{F: intPtr(zeroInt)})},
		{"int", true, false, "no tag", "non-zero", marshal(IntPtrNoTag{F: intPtr(nonZeroInt)})},
		{"int", true, false, "tag", "nil", marshal(IntPtrTag{})},
		{"int", true, false, "tag", "zero", marshal(IntPtrTag{F: intPtr(zeroInt)})},
		{"int", true, false, "tag", "non-zero", marshal(IntPtrTag{F: intPtr(nonZeroInt)})},
		{"int", true, true, "omitempty", "nil", marshal(IntPtrOmit{})},
		{"int", true, true, "omitempty", "zero", marshal(IntPtrOmit{F: intPtr(zeroInt)})},
		{"int", true, true, "omitempty", "non-zero", marshal(IntPtrOmit{F: intPtr(nonZeroInt)})},

		// ── bool ───────────────────────────────────────────────────────────
		{"bool", false, false, "no tag", "zero", marshal(BoolNoTag{})},
		{"bool", false, false, "no tag", "non-zero", marshal(BoolNoTag{F: nonZeroBool})},
		{"bool", false, false, "tag", "zero", marshal(BoolTag{})},
		{"bool", false, false, "tag", "non-zero", marshal(BoolTag{F: nonZeroBool})},
		{"bool", false, true, "omitempty", "zero", marshal(BoolOmit{})},
		{"bool", false, true, "omitempty", "non-zero", marshal(BoolOmit{F: nonZeroBool})},
		{"bool", true, false, "no tag", "nil", marshal(BoolPtrNoTag{})},
		{"bool", true, false, "no tag", "zero", marshal(BoolPtrNoTag{F: boolPtr(zeroBool)})},
		{"bool", true, false, "no tag", "non-zero", marshal(BoolPtrNoTag{F: boolPtr(nonZeroBool)})},
		{"bool", true, false, "tag", "nil", marshal(BoolPtrTag{})},
		{"bool", true, false, "tag", "zero", marshal(BoolPtrTag{F: boolPtr(zeroBool)})},
		{"bool", true, false, "tag", "non-zero", marshal(BoolPtrTag{F: boolPtr(nonZeroBool)})},
		{"bool", true, true, "omitempty", "nil", marshal(BoolPtrOmit{})},
		{"bool", true, true, "omitempty", "zero", marshal(BoolPtrOmit{F: boolPtr(zeroBool)})},
		{"bool", true, true, "omitempty", "non-zero", marshal(BoolPtrOmit{F: boolPtr(nonZeroBool)})},

		// ── float64 ────────────────────────────────────────────────────────
		{"float64", false, false, "no tag", "zero", marshal(Float64NoTag{})},
		{"float64", false, false, "no tag", "non-zero", marshal(Float64NoTag{F: nonZeroFloat})},
		{"float64", false, false, "tag", "zero", marshal(Float64Tag{})},
		{"float64", false, false, "tag", "non-zero", marshal(Float64Tag{F: nonZeroFloat})},
		{"float64", false, true, "omitempty", "zero", marshal(Float64Omit{})},
		{"float64", false, true, "omitempty", "non-zero", marshal(Float64Omit{F: nonZeroFloat})},
		{"float64", true, false, "no tag", "nil", marshal(Float64PtrNoTag{})},
		{"float64", true, false, "no tag", "zero", marshal(Float64PtrNoTag{F: float64Ptr(zeroFloat)})},
		{"float64", true, false, "no tag", "non-zero", marshal(Float64PtrNoTag{F: float64Ptr(nonZeroFloat)})},
		{"float64", true, false, "tag", "nil", marshal(Float64PtrTag{})},
		{"float64", true, false, "tag", "zero", marshal(Float64PtrTag{F: float64Ptr(zeroFloat)})},
		{"float64", true, false, "tag", "non-zero", marshal(Float64PtrTag{F: float64Ptr(nonZeroFloat)})},
		{"float64", true, true, "omitempty", "nil", marshal(Float64PtrOmit{})},
		{"float64", true, true, "omitempty", "zero", marshal(Float64PtrOmit{F: float64Ptr(zeroFloat)})},
		{"float64", true, true, "omitempty", "non-zero", marshal(Float64PtrOmit{F: float64Ptr(nonZeroFloat)})},

		// ── []string (slice) ───────────────────────────────────────────────
		// Note: nil slice and empty slice differ with omitempty
		{"[]string", false, false, "no tag", "nil", marshal(SliceNoTag{F: nilSlice})},
		{"[]string", false, false, "no tag", "empty", marshal(SliceNoTag{F: zeroSlice})},
		{"[]string", false, false, "no tag", "non-empty", marshal(SliceNoTag{F: nonZeroSlice})},
		{"[]string", false, false, "tag", "nil", marshal(SliceTag{F: nilSlice})},
		{"[]string", false, false, "tag", "empty", marshal(SliceTag{F: zeroSlice})},
		{"[]string", false, false, "tag", "non-empty", marshal(SliceTag{F: nonZeroSlice})},
		{"[]string", false, true, "omitempty", "nil", marshal(SliceOmit{F: nilSlice})},
		{"[]string", false, true, "omitempty", "empty", marshal(SliceOmit{F: zeroSlice})},
		{"[]string", false, true, "omitempty", "non-empty", marshal(SliceOmit{F: nonZeroSlice})},
		{"[]string", true, false, "no tag", "nil ptr", marshal(SlicePtrNoTag{})},
		{"[]string", true, false, "no tag", "ptr→nil", marshal(SlicePtrNoTag{F: slicePtr(nilSlice)})},
		{"[]string", true, false, "no tag", "ptr→empty", marshal(SlicePtrNoTag{F: slicePtr(zeroSlice)})},
		{"[]string", true, false, "no tag", "ptr→non-empty", marshal(SlicePtrNoTag{F: slicePtr(nonZeroSlice)})},
		{"[]string", true, false, "tag", "nil ptr", marshal(SlicePtrTag{})},
		{"[]string", true, false, "tag", "ptr→empty", marshal(SlicePtrTag{F: slicePtr(zeroSlice)})},
		{"[]string", true, false, "tag", "ptr→non-empty", marshal(SlicePtrTag{F: slicePtr(nonZeroSlice)})},
		{"[]string", true, true, "omitempty", "nil ptr", marshal(SlicePtrOmit{})},
		{"[]string", true, true, "omitempty", "ptr→nil", marshal(SlicePtrOmit{F: slicePtr(nilSlice)})},
		{"[]string", true, true, "omitempty", "ptr→empty", marshal(SlicePtrOmit{F: slicePtr(zeroSlice)})},
		{"[]string", true, true, "omitempty", "ptr→non-empty", marshal(SlicePtrOmit{F: slicePtr(nonZeroSlice)})},

		// ── map[string]string ──────────────────────────────────────────────
		{"map", false, false, "no tag", "nil", marshal(MapNoTag{})},
		{"map", false, false, "no tag", "empty", marshal(MapNoTag{F: zeroMap})},
		{"map", false, false, "no tag", "non-empty", marshal(MapNoTag{F: nonZeroMap})},
		{"map", false, false, "tag", "nil", marshal(MapTag{})},
		{"map", false, false, "tag", "empty", marshal(MapTag{F: zeroMap})},
		{"map", false, false, "tag", "non-empty", marshal(MapTag{F: nonZeroMap})},
		{"map", false, true, "omitempty", "nil", marshal(MapOmit{})},
		{"map", false, true, "omitempty", "empty", marshal(MapOmit{F: zeroMap})},
		{"map", false, true, "omitempty", "non-empty", marshal(MapOmit{F: nonZeroMap})},
		{"map", true, false, "no tag", "nil ptr", marshal(MapPtrNoTag{})},
		{"map", true, false, "no tag", "ptr→nil", marshal(MapPtrNoTag{F: mapPtr(nil)})},
		{"map", true, false, "no tag", "ptr→empty", marshal(MapPtrNoTag{F: mapPtr(zeroMap)})},
		{"map", true, false, "no tag", "ptr→non-empty", marshal(MapPtrNoTag{F: mapPtr(nonZeroMap)})},
		{"map", true, false, "tag", "nil ptr", marshal(MapPtrTag{})},
		{"map", true, false, "tag", "ptr→empty", marshal(MapPtrTag{F: mapPtr(zeroMap)})},
		{"map", true, false, "tag", "ptr→non-empty", marshal(MapPtrTag{F: mapPtr(nonZeroMap)})},
		{"map", true, true, "omitempty", "nil ptr", marshal(MapPtrOmit{})},
		{"map", true, true, "omitempty", "ptr→nil", marshal(MapPtrOmit{F: mapPtr(nil)})},
		{"map", true, true, "omitempty", "ptr→empty", marshal(MapPtrOmit{F: mapPtr(zeroMap)})},
		{"map", true, true, "omitempty", "ptr→non-empty", marshal(MapPtrOmit{F: mapPtr(nonZeroMap)})},

		// ── struct ─────────────────────────────────────────────────────────
		// Note: omitempty does NOT omit a zero struct (it is never "empty")
		{"struct", false, false, "no tag", "zero", marshal(StructNoTag{})},
		{"struct", false, false, "no tag", "non-zero", marshal(StructNoTag{F: nonZeroInner})},
		{"struct", false, false, "tag", "zero", marshal(StructTag{})},
		{"struct", false, false, "tag", "non-zero", marshal(StructTag{F: nonZeroInner})},
		{"struct", false, true, "omitempty", "zero", marshal(StructOmit{})},      // ← NOT omitted!
		{"struct", false, true, "omitempty", "non-zero", marshal(StructOmit{F: nonZeroInner})},
		{"struct", true, false, "no tag", "nil", marshal(StructPtrNoTag{})},
		{"struct", true, false, "no tag", "zero", marshal(StructPtrNoTag{F: innerPtr(zeroInner)})},
		{"struct", true, false, "no tag", "non-zero", marshal(StructPtrNoTag{F: innerPtr(nonZeroInner)})},
		{"struct", true, false, "tag", "nil", marshal(StructPtrTag{})},
		{"struct", true, false, "tag", "zero", marshal(StructPtrTag{F: innerPtr(zeroInner)})},
		{"struct", true, false, "tag", "non-zero", marshal(StructPtrTag{F: innerPtr(nonZeroInner)})},
		{"struct", true, true, "omitempty", "nil", marshal(StructPtrOmit{})},
		{"struct", true, true, "omitempty", "zero", marshal(StructPtrOmit{F: innerPtr(zeroInner)})},
		{"struct", true, true, "omitempty", "non-zero", marshal(StructPtrOmit{F: innerPtr(nonZeroInner)})},
	}

	// Print as a markdown table for easy reading in test output
	header := fmt.Sprintf("%-12s %-8s %-12s %-14s %s",
		"type", "pointer", "tag", "value", "JSON output")
	sep := strings.Repeat("-", len(header)+20)
	t.Log("\n" + sep)
	t.Log(header)
	t.Log(sep)
	for _, r := range rows {
		ptrStr := "no"
		if r.ptr {
			ptrStr = "yes"
		}
		t.Logf("%-12s %-8s %-12s %-14s %s",
			r.label, ptrStr, r.tagName, r.valueKind, r.result)
	}
	t.Log(sep)
}
