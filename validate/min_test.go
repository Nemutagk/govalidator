package validate

import "testing"

func TestMin_InvalidOption_NotANumber(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Min("qty", 5, map[string]any{}, []string{"abc"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "qty", "min")
	want := "El campo qty debe ser un número"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestMin_InvalidOption_NotANumber_WithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Min("qty", 5, map[string]any{}, []string{"abc"}, "1", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "qty", "min")
	want := "El campo qty en la posición 1 debe ser un número"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestMin_InvalidOption_NotANumber_CustomError(t *testing.T) {
	errors := make(map[string]interface{})
	customErrors := map[string]string{"qty.min": "mensaje personalizado"}
	errors = Min("qty", 5, map[string]any{}, []string{"abc"}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "qty", "min")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestMin_String_Passes(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Min("name", "hello", map[string]any{}, []string{"3"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestMin_String_EqualLengthPasses(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Min("name", "abc", map[string]any{}, []string{"3"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestMin_String_Fails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Min("name", "hi", map[string]any{}, []string{"3"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "min")
	want := "El campo name debe tener al menos 3 caracteres"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestMin_String_Fails_WithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Min("name", "hi", map[string]any{}, []string{"3"}, "2", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "name", "min")
	want := "El campo name en la posición 2 debe tener al menos 3 caracteres"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestMin_String_Fails_CustomError(t *testing.T) {
	errors := make(map[string]interface{})
	customErrors := map[string]string{"name.min": "mensaje personalizado"}
	errors = Min("name", "hi", map[string]any{}, []string{"3"}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "name", "min")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestMin_Int_Passes(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Min("age", 20, map[string]any{}, []string{"18"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestMin_Int_Fails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Min("age", 10, map[string]any{}, []string{"18"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "age", "min")
	want := "El campo age debe ser al menos 18"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestMin_Int_Fails_WithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Min("age", 10, map[string]any{}, []string{"18"}, "4", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "age", "min")
	want := "El campo age en la posición 4 debe ser al menos 18"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestMin_Int_Fails_CustomError(t *testing.T) {
	errors := make(map[string]interface{})
	customErrors := map[string]string{"age.min": "mensaje personalizado"}
	errors = Min("age", 10, map[string]any{}, []string{"18"}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "age", "min")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestMin_Float64_Passes(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Min("price", 15.5, map[string]any{}, []string{"10"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestMin_Float64_Fails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Min("price", 5.5, map[string]any{}, []string{"10"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "price", "min")
	want := "El campo price debe ser al menos 10"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestMin_Float64_Fails_WithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Min("price", 5.5, map[string]any{}, []string{"10"}, "3", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "price", "min")
	want := "El campo price en la posición 3 debe ser al menos 10"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestMin_Float64_Fails_CustomError(t *testing.T) {
	errors := make(map[string]interface{})
	customErrors := map[string]string{"price.min": "mensaje personalizado"}
	errors = Min("price", 5.5, map[string]any{}, []string{"10"}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "price", "min")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestMin_Slice_Passes(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Min("items", []string{"a", "b", "c"}, map[string]any{}, []string{"3"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestMin_Slice_Fails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Min("items", []string{"a"}, map[string]any{}, []string{"3"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "items", "min")
	want := "El campo items debe tener al menos 3 elementos"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestMin_Slice_Fails_WithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Min("items", []string{"a"}, map[string]any{}, []string{"3"}, "0", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "items", "min")
	want := "El campo items en la posición 0 debe tener al menos 3 elementos"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestMin_Slice_Fails_CustomError(t *testing.T) {
	errors := make(map[string]interface{})
	customErrors := map[string]string{"items.min": "mensaje personalizado"}
	errors = Min("items", []string{"a"}, map[string]any{}, []string{"3"}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "items", "min")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestMin_Array_Fails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Min("items", [2]int{1, 2}, map[string]any{}, []string{"3"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "items", "min")
	want := "El campo items debe tener al menos 3 elementos"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestMin_Array_Passes(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Min("items", [3]int{1, 2, 3}, map[string]any{}, []string{"3"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestMin_UnsupportedType_NoError(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Min("active", true, map[string]any{}, []string{"3"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors for unsupported type, got %v", errors)
	}
}

func TestMin_NilValue_NoError(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Min("field", nil, map[string]any{}, []string{"3"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors for nil value, got %v", errors)
	}
}
