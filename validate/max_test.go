package validate

import "testing"

func TestMax_InvalidOption_NotANumber(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Max("qty", 5, map[string]any{}, []string{"abc"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "qty", "max")
	want := "El campo qty debe ser un número"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestMax_InvalidOption_NotANumber_WithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Max("qty", 5, map[string]any{}, []string{"abc"}, "1", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "qty", "max")
	want := "El campo qty en la posición 1 debe ser un número"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestMax_InvalidOption_NotANumber_CustomError(t *testing.T) {
	errors := make(map[string]interface{})
	customErrors := map[string]string{"qty.max": "mensaje personalizado"}
	errors = Max("qty", 5, map[string]any{}, []string{"abc"}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "qty", "max")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestMax_String_Passes(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Max("name", "hi", map[string]any{}, []string{"5"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestMax_String_EqualLengthPasses(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Max("name", "hello", map[string]any{}, []string{"5"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestMax_String_Fails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Max("name", "helloooo", map[string]any{}, []string{"5"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "max")
	want := "El campo name debe tener como máximo 5 caracteres"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestMax_String_Fails_WithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Max("name", "helloooo", map[string]any{}, []string{"5"}, "2", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "max")
	want := "El campo name en la posición 2 debe tener como máximo 5 caracteres"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestMax_String_Fails_CustomError(t *testing.T) {
	errors := make(map[string]interface{})
	customErrors := map[string]string{"name.max": "mensaje personalizado"}
	errors = Max("name", "helloooo", map[string]any{}, []string{"5"}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "name", "max")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestMax_Int_Passes(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Max("age", 40, map[string]any{}, []string{"65"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestMax_Int_Fails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Max("age", 70, map[string]any{}, []string{"65"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "age", "max")
	want := "El campo age debe ser como máximo 65"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestMax_Int_Fails_WithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Max("age", 70, map[string]any{}, []string{"65"}, "4", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "age", "max")
	want := "El campo age en la posición 4 debe ser como máximo 65"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestMax_Int_Fails_CustomError(t *testing.T) {
	errors := make(map[string]interface{})
	customErrors := map[string]string{"age.max": "mensaje personalizado"}
	errors = Max("age", 70, map[string]any{}, []string{"65"}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "age", "max")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestMax_Float64_Passes(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Max("price", 50.5, map[string]any{}, []string{"100"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestMax_Float64_Fails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Max("price", 150.5, map[string]any{}, []string{"100"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "price", "max")
	want := "El campo price debe ser como máximo 100"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestMax_Float64_Fails_WithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Max("price", 150.5, map[string]any{}, []string{"100"}, "3", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "price", "max")
	want := "El campo price en la posición 3 debe ser como máximo 100"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestMax_Float64_Fails_CustomError(t *testing.T) {
	errors := make(map[string]interface{})
	customErrors := map[string]string{"price.max": "mensaje personalizado"}
	errors = Max("price", 150.5, map[string]any{}, []string{"100"}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "price", "max")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestMax_Slice_Passes(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Max("items", []string{"a"}, map[string]any{}, []string{"2"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestMax_Slice_EqualLengthPasses(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Max("items", []string{"a", "b"}, map[string]any{}, []string{"2"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestMax_Slice_Fails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Max("items", []string{"a", "b", "c"}, map[string]any{}, []string{"2"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "items", "max")
	want := "El campo items debe tener como máximo 2 elementos"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestMax_Slice_Fails_WithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Max("items", []string{"a", "b", "c"}, map[string]any{}, []string{"2"}, "0", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "items", "max")
	want := "El campo items en la posición 0 debe tener como máximo 2 elementos"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestMax_Slice_Fails_CustomError(t *testing.T) {
	errors := make(map[string]interface{})
	customErrors := map[string]string{"items.max": "mensaje personalizado"}
	errors = Max("items", []string{"a", "b", "c"}, map[string]any{}, []string{"2"}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "items", "max")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestMax_Array_Fails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Max("items", [3]int{1, 2, 3}, map[string]any{}, []string{"2"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "items", "max")
	want := "El campo items debe tener como máximo 2 elementos"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestMax_Array_Passes(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Max("items", [2]int{1, 2}, map[string]any{}, []string{"2"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestMax_UnsupportedType_NoError(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Max("active", true, map[string]any{}, []string{"2"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors for unsupported type, got %v", errors)
	}
}

func TestMax_NilValue_NoError(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Max("field", nil, map[string]any{}, []string{"2"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors for nil value, got %v", errors)
	}
}
