package validate

import (
	"fmt"
	"math"
	"testing"
)

type namedCount uint16

type numericRule func(input string, value any, payload map[string]any, options []string, sliceIndex string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, customeErrors map[string]string) map[string]interface{}

// numericSamples devuelve el valor 5 representado en cada tipo numérico de Go.
func numericSamples() []any {
	return []any{
		int(5), int8(5), int16(5), int32(5), int64(5),
		uint(5), uint8(5), uint16(5), uint32(5), uint64(5),
		float32(5), float64(5), namedCount(5),
	}
}

func TestNumericRules_AllNumericTypes(t *testing.T) {
	cases := []struct {
		rule    string
		fn      numericRule
		option  string
		wantErr bool
	}{
		{"min", Min, "1", false},
		{"min", Min, "5", false},
		{"min", Min, "10", true},
		{"max", Max, "10", false},
		{"max", Max, "5", false},
		{"max", Max, "3", true},
		{"greater_than", GreaterThan, "3", false},
		{"greater_than", GreaterThan, "5", true},
		{"greater_than_equal", GreaterThanEqual, "5", false},
		{"greater_than_equal", GreaterThanEqual, "6", true},
		{"less_than", LessThan, "10", false},
		{"less_than", LessThan, "5", true},
		{"less_than_equal", LessThanEqual, "5", false},
		{"less_than_equal", LessThanEqual, "4", true},
		{"before", Before, "10", false},
		{"before", Before, "3", true},
		{"after", After, "3", false},
		{"after", After, "10", true},
	}

	for _, c := range cases {
		for _, value := range numericSamples() {
			t.Run(fmt.Sprintf("%s_%s_%T", c.rule, c.option, value), func(t *testing.T) {
				errors := make(map[string]interface{})
				errors = c.fn("n", value, map[string]any{}, []string{c.option}, "", errors, testAddError, map[string]string{})

				_, found := getErrorMsgs(errors, "n", c.rule)
				if found != c.wantErr {
					t.Fatalf("value %v (%T) rule %s:%s: error = %v, want %v (errors: %v)", value, value, c.rule, c.option, found, c.wantErr, errors)
				}
			})
		}
	}
}

func TestMin_Int64BelowMinFails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Min("kdf_memory_kib", int64(1024), map[string]any{}, []string{"19456"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "kdf_memory_kib", "min")
	want := "El campo kdf_memory_kib debe ser al menos 19456"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestMin_Max_Int64_WithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Min("qty", int64(1), map[string]any{}, []string{"2"}, "3", errors, testAddError, map[string]string{})
	msgs, found := getErrorMsgs(errors, "qty", "min")
	want := "El campo qty en la posición 3 debe ser al menos 2"
	if !found || msgs[0] != want {
		t.Fatalf("min msgs = %v, want %q", msgs, want)
	}

	errors = make(map[string]interface{})
	errors = Max("qty", int64(9), map[string]any{}, []string{"2"}, "3", errors, testAddError, map[string]string{})
	msgs, found = getErrorMsgs(errors, "qty", "max")
	want = "El campo qty en la posición 3 debe ser como máximo 2"
	if !found || msgs[0] != want {
		t.Fatalf("max msgs = %v, want %q", msgs, want)
	}
}

func TestMin_Max_Int64_CustomError(t *testing.T) {
	customErrors := map[string]string{"qty.min": "min personalizado", "qty.max": "max personalizado"}

	errors := make(map[string]interface{})
	errors = Min("qty", int64(1), map[string]any{}, []string{"2"}, "", errors, testAddError, customErrors)
	if msgs, found := getErrorMsgs(errors, "qty", "min"); !found || msgs[0] != "min personalizado" {
		t.Fatalf("expected custom min message, got %v", errors)
	}

	errors = make(map[string]interface{})
	errors = Max("qty", int64(9), map[string]any{}, []string{"2"}, "", errors, testAddError, customErrors)
	if msgs, found := getErrorMsgs(errors, "qty", "max"); !found || msgs[0] != "max personalizado" {
		t.Fatalf("expected custom max message, got %v", errors)
	}
}

func TestMinMax_Uint64BeyondInt64Range(t *testing.T) {
	big := uint64(math.MaxUint64)

	errors := make(map[string]interface{})
	errors = Min("n", big, map[string]any{}, []string{"100"}, "", errors, testAddError, map[string]string{})
	if len(errors) != 0 {
		t.Fatalf("min: expected no errors for MaxUint64, got %v", errors)
	}

	errors = make(map[string]interface{})
	errors = Max("n", big, map[string]any{}, []string{"100"}, "", errors, testAddError, map[string]string{})
	if _, found := getErrorMsgs(errors, "n", "max"); !found {
		t.Fatalf("max: expected an error for MaxUint64, got %v", errors)
	}
}

func TestMinMax_Int64ExtremesAreExact(t *testing.T) {
	// float64 no distingue math.MaxInt64 de math.MaxInt64-1; la comparación debe ser entera.
	errors := make(map[string]interface{})
	errors = Max("n", int64(math.MaxInt64), map[string]any{}, []string{"9223372036854775806"}, "", errors, testAddError, map[string]string{})
	if _, found := getErrorMsgs(errors, "n", "max"); !found {
		t.Fatalf("expected an error for MaxInt64 over a max of MaxInt64-1, got %v", errors)
	}
}

func TestMinMax_Float32KeepsDecimalValue(t *testing.T) {
	errors := make(map[string]interface{})
	errors = LessThanEqual("n", float32(0.1), map[string]any{}, []string{"0.1"}, "", errors, testAddError, map[string]string{})
	if len(errors) != 0 {
		t.Fatalf("expected float32(0.1) <= 0.1 to pass, got %v", errors)
	}

	errors = make(map[string]interface{})
	errors = GreaterThanEqual("n", float32(0.1), map[string]any{}, []string{"0.1"}, "", errors, testAddError, map[string]string{})
	if len(errors) != 0 {
		t.Fatalf("expected float32(0.1) >= 0.1 to pass, got %v", errors)
	}
}

func TestMinMax_NonNumericAndNilAreIgnored(t *testing.T) {
	for _, v := range []any{nil, true, map[string]any{"k": 1}} {
		errors := make(map[string]interface{})
		errors = Min("n", v, map[string]any{}, []string{"2"}, "", errors, testAddError, map[string]string{})
		errors = Max("n", v, map[string]any{}, []string{"2"}, "", errors, testAddError, map[string]string{})
		if len(errors) != 0 {
			t.Fatalf("value %v: expected no errors, got %v", v, errors)
		}
	}
}

func TestGreaterThan_NumericFieldTargetOfAnotherType(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"max_qty": uint8(3)}
	errors = GreaterThan("qty", int32(5), payload, []string{"max_qty"}, "", errors, testAddError, map[string]string{})
	if len(errors) != 0 {
		t.Fatalf("expected int32(5) > uint8(3) to pass, got %v", errors)
	}
}

func TestToFloat64_AllNumericTypes(t *testing.T) {
	for _, v := range numericSamples() {
		got, ok := toFloat64(v)
		if !ok || got != 5 {
			t.Fatalf("toFloat64(%v (%T)) = (%v, %v), want (5, true)", v, v, got, ok)
		}
	}
}

func TestCompareToInt64(t *testing.T) {
	if _, ok := compareToInt64("5", 1); ok {
		t.Fatalf("string must not be numeric")
	}
	if _, ok := compareToInt64(nil, 1); ok {
		t.Fatalf("nil must not be numeric")
	}
	if _, ok := compareToInt64(true, 1); ok {
		t.Fatalf("bool must not be numeric")
	}

	cases := []struct {
		value any
		limit int64
		want  int
	}{
		{int8(-3), -3, 0},
		{int8(-4), -3, -1},
		{uint8(200), 100, 1},
		{uint64(math.MaxUint64), math.MaxInt64, 1},
		{float64(2.5), 2, 1},
		{float64(1.5), 2, -1},
	}
	for _, c := range cases {
		got, ok := compareToInt64(c.value, c.limit)
		if !ok || got != c.want {
			t.Fatalf("compareToInt64(%v (%T), %d) = (%d, %v), want (%d, true)", c.value, c.value, c.limit, got, ok, c.want)
		}
	}
}
