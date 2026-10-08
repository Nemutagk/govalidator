package validate

import (
	"fmt"
	"testing"
)

type namedLabel string

func TestEqualsOption(t *testing.T) {
	cases := []struct {
		value  any
		option string
		want   bool
	}{
		{"abc", "abc", true},
		{"abc", "abd", false},
		{namedLabel("abc"), "abc", true},
		{5, "5", true},
		{int8(-5), "-5", true},
		{uint64(18446744073709551615), "18446744073709551615", true},
		{float64(5), "5", true},
		{5.5, "5.5", true},
		{float32(0.1), "0.1", true},
		{true, "true", true},
		{false, "true", false},
		{5, "6", false},
		{5, "5.0", false},
		{nil, "", false},
		{[]string{"a"}, "a", false},
		{map[string]any{"a": 1}, "a", false},
	}

	for _, c := range cases {
		if got := equalsOption(c.value, c.option); got != c.want {
			t.Fatalf("equalsOption(%v (%T), %q) = %v, want %v", c.value, c.value, c.option, got, c.want)
		}
	}
}

func TestIn_NonStringValues_AreComparedAsText(t *testing.T) {
	for _, v := range []any{5, int64(5), uint8(5), float64(5), float32(5)} {
		errors := make(map[string]interface{})
		errors = In("n", v, map[string]any{}, []string{"4", "5"}, "", errors, testAddError, map[string]string{})
		if len(errors) != 0 {
			t.Fatalf("value %v (%T): expected no errors, got %v", v, v, errors)
		}
	}

	errors := make(map[string]interface{})
	errors = In("n", int64(9), map[string]any{}, []string{"4", "5"}, "", errors, testAddError, map[string]string{})
	if _, found := getErrorMsgs(errors, "n", "in"); !found {
		t.Fatalf("expected an error for a value outside the options, got %v", errors)
	}
}

func TestIn_Bool_IsComparedAsText(t *testing.T) {
	errors := make(map[string]interface{})
	errors = In("flag", true, map[string]any{}, []string{"true", "false"}, "", errors, testAddError, map[string]string{})
	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestNotIn_NonStringValues_AreComparedAsText(t *testing.T) {
	errors := make(map[string]interface{})
	errors = NotIn("n", int64(0), map[string]any{}, []string{"0"}, "", errors, testAddError, map[string]string{})
	if _, found := getErrorMsgs(errors, "n", "notin"); !found {
		t.Fatalf("expected an error for a forbidden numeric value, got %v", errors)
	}

	errors = make(map[string]interface{})
	errors = NotIn("n", int64(1), map[string]any{}, []string{"0"}, "", errors, testAddError, map[string]string{})
	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestEqual_NonStringValues_AreComparedAsText(t *testing.T) {
	for _, v := range []any{5, int64(5), uint16(5), float64(5), true} {
		option := fmt.Sprint(v)
		errors := make(map[string]interface{})
		errors = Equal("n", v, map[string]any{"n": v}, []string{option}, "", errors, testAddError, map[string]string{})
		if len(errors) != 0 {
			t.Fatalf("value %v (%T): expected no errors, got %v", v, v, errors)
		}
	}

	errors := make(map[string]interface{})
	errors = Equal("n", int64(5), map[string]any{"n": int64(5)}, []string{"6"}, "", errors, testAddError, map[string]string{})
	if _, found := getErrorMsgs(errors, "n", "equal"); !found {
		t.Fatalf("expected an error when the value differs, got %v", errors)
	}
}

func TestNotEqual_NonStringValues_AreComparedAsText(t *testing.T) {
	errors := make(map[string]interface{})
	errors = NotEqual("n", int64(0), map[string]any{"n": int64(0)}, []string{"0"}, "", errors, testAddError, map[string]string{})
	if _, found := getErrorMsgs(errors, "n", "not_equal"); !found {
		t.Fatalf("expected an error for a value equal to the option, got %v", errors)
	}

	errors = make(map[string]interface{})
	errors = NotEqual("n", int64(1), map[string]any{"n": int64(1)}, []string{"0"}, "", errors, testAddError, map[string]string{})
	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestInIf_NonStringValues_AreComparedAsText(t *testing.T) {
	body := map[string]any{"mode": "a"}

	errors := make(map[string]interface{})
	errors = InIf("kind", int64(3), body, []string{"mode", "a", "3", "4"}, "", errors, testAddError, map[string]string{})
	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}

	errors = make(map[string]interface{})
	errors = InIf("kind", int64(9), body, []string{"mode", "a", "3", "4"}, "", errors, testAddError, map[string]string{})
	if _, found := getErrorMsgs(errors, "kind", "in_if"); !found {
		t.Fatalf("expected an error for a value outside the allowed list, got %v", errors)
	}
}

func TestRequiredWith_ExpectedValueOnNumericField(t *testing.T) {
	errors := make(map[string]interface{})
	errors, failed := RequiredWith("detail", nil, map[string]any{"level": 5}, []string{"level", "5"}, "", errors, testAddError, map[string]string{})
	if _, found := getErrorMsgs(errors, "detail", "required_with"); !found || !failed {
		t.Fatalf("expected detail to be required when level is 5, got %v", errors)
	}

	errors = make(map[string]interface{})
	errors, failed = RequiredWith("detail", nil, map[string]any{"level": 4}, []string{"level", "5"}, "", errors, testAddError, map[string]string{})
	if len(errors) != 0 || failed {
		t.Fatalf("expected detail to be optional when level is 4, got %v", errors)
	}
}

func TestNotIn_CustomErrorAcceptsBothKeys(t *testing.T) {
	for _, key := range []string{"status.notin", "status.not_in"} {
		errors := make(map[string]interface{})
		errors = NotIn("status", "x", map[string]any{}, []string{"x"}, "", errors, testAddError, map[string]string{key: "personalizado"})
		if msgs, found := getErrorMsgs(errors, "status", "notin"); !found || msgs[0] != "personalizado" {
			t.Fatalf("key %s: expected the custom message, got %v", key, errors)
		}
	}
}

func TestNullable_CustomErrorAcceptsBothKeys(t *testing.T) {
	for _, key := range []string{"name.null", "name.nullable"} {
		errors := make(map[string]interface{})
		errors = Nullable("name", nil, map[string]any{}, []string{}, "", errors, testAddError, map[string]string{key: "personalizado"})
		if msgs, found := getErrorMsgs(errors, "name", "null"); !found || msgs[0] != "personalizado" {
			t.Fatalf("key %s: expected the custom message, got %v", key, errors)
		}
	}
}

func TestBeforeAfter_NonNumericOptionOnNumber(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Before("n", 5, map[string]any{}, []string{"abc"}, "", errors, testAddError, map[string]string{})
	msgs, found := getErrorMsgs(errors, "n", "before")
	want := "El valor a comparar abc no es un número válido"
	if !found || msgs[0] != want {
		t.Fatalf("before msgs = %v, want %q", msgs, want)
	}

	errors = make(map[string]interface{})
	errors = After("n", 5, map[string]any{}, []string{"abc"}, "", errors, testAddError, map[string]string{})
	msgs, found = getErrorMsgs(errors, "n", "after")
	if !found || msgs[0] != want {
		t.Fatalf("after msgs = %v, want %q", msgs, want)
	}
}

func TestBeforeAfter_NumericFieldAsTarget(t *testing.T) {
	payload := map[string]any{"limit": int64(10)}

	errors := make(map[string]interface{})
	errors = Before("n", 5, payload, []string{"limit"}, "", errors, testAddError, map[string]string{})
	if len(errors) != 0 {
		t.Fatalf("before: expected 5 < limit(10) to pass, got %v", errors)
	}

	errors = make(map[string]interface{})
	errors = Before("n", 15, payload, []string{"limit"}, "", errors, testAddError, map[string]string{})
	msgs, found := getErrorMsgs(errors, "n", "before")
	want := "La entrada n no es anterior al campo limit"
	if !found || msgs[0] != want {
		t.Fatalf("before msgs = %v, want %q", msgs, want)
	}

	errors = make(map[string]interface{})
	errors = After("n", 15, payload, []string{"limit"}, "", errors, testAddError, map[string]string{})
	if len(errors) != 0 {
		t.Fatalf("after: expected 15 > limit(10) to pass, got %v", errors)
	}

	errors = make(map[string]interface{})
	errors = After("n", 5, payload, []string{"limit"}, "2", errors, testAddError, map[string]string{})
	msgs, found = getErrorMsgs(errors, "n", "after")
	want = "La entrada en la posición 2 no es posterior al campo limit"
	if !found || msgs[0] != want {
		t.Fatalf("after msgs = %v, want %q", msgs, want)
	}
}

func TestBeforeAfter_NonNumericFieldAsTarget(t *testing.T) {
	payload := map[string]any{"limit": "diez"}

	errors := make(map[string]interface{})
	errors = After("n", 5, payload, []string{"limit"}, "", errors, testAddError, map[string]string{})
	msgs, found := getErrorMsgs(errors, "n", "after")
	want := "El campo limit no es un número válido para comparar"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestBeforeAfter_DecimalLiteralOnNumber(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Before("n", 2.4, map[string]any{}, []string{"2.5"}, "", errors, testAddError, map[string]string{})
	if len(errors) != 0 {
		t.Fatalf("expected 2.4 < 2.5 to pass, got %v", errors)
	}

	errors = make(map[string]interface{})
	errors = After("n", 2.4, map[string]any{}, []string{"2.5"}, "", errors, testAddError, map[string]string{})
	if _, found := getErrorMsgs(errors, "n", "after"); !found {
		t.Fatalf("expected 2.4 > 2.5 to fail, got %v", errors)
	}
}
