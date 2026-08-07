package validate

import "testing"

func TestArray_SlicePasses(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Array("items", []string{"a", "b"}, map[string]any{}, []string{}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestArray_FixedArrayPasses(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Array("items", [2]int{1, 2}, map[string]any{}, []string{}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestArray_EmptySlicePasses(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Array("items", []string{}, map[string]any{}, []string{}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestArray_String_Fails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Array("items", "hello", map[string]any{}, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "items", "array")
	want := "El valor no es un array."
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestArray_String_Fails_WithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Array("items", "hello", map[string]any{}, []string{}, "2", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "items", "array")
	want := "El campo items en la posición 2 no es un array."
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestArray_String_Fails_CustomError(t *testing.T) {
	errors := make(map[string]interface{})
	customErrors := map[string]string{"items.array": "mensaje personalizado"}
	errors = Array("items", "hello", map[string]any{}, []string{}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "items", "array")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestArray_Int_Fails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Array("items", 5, map[string]any{}, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "items", "array")
	want := "El valor no es un array."
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestArray_Map_Fails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Array("items", map[string]any{"a": 1}, map[string]any{}, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "items", "array")
	want := "El valor no es un array."
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestArray_Nil_Fails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Array("items", nil, map[string]any{}, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "items", "array")
	want := "El valor no es un array."
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}
