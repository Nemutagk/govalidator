package validate

import "testing"

func TestLen_MissingOption(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Len("code", "abc", map[string]any{}, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "code", "len")
	want := "La regla len del campo code requiere un valor"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestLen_MissingOption_WithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Len("code", "abc", map[string]any{}, nil, "2", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "code", "len")
	want := "La regla len del campo code en la posición 2 requiere un valor"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestLen_MissingOption_CustomError(t *testing.T) {
	errors := make(map[string]interface{})
	customErrors := map[string]string{"code.len": "mensaje personalizado"}
	errors = Len("code", "abc", map[string]any{}, []string{}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "code", "len")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestLen_InvalidOption_NotANumber(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Len("code", "abc", map[string]any{}, []string{"abc"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "code", "len")
	want := "El campo code debe ser un número"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestLen_InvalidOption_NotANumber_WithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Len("code", "abc", map[string]any{}, []string{"abc"}, "1", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "code", "len")
	want := "El campo code en la posición 1 debe ser un número"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestLen_InvalidOption_NotANumber_CustomError(t *testing.T) {
	errors := make(map[string]interface{})
	customErrors := map[string]string{"code.len": "mensaje personalizado"}
	errors = Len("code", "abc", map[string]any{}, []string{"abc"}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "code", "len")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestLen_String_Passes(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Len("code", "abc", map[string]any{}, []string{"3"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestLen_String_CountsCharactersNotBytes(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Len("name", "añoñ", map[string]any{}, []string{"4"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestLen_String_EmptyPassesWithZero(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Len("code", "", map[string]any{}, []string{"0"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestLen_String_FailsWhenShorter(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Len("code", "ab", map[string]any{}, []string{"3"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "code", "len")
	want := "El campo code debe tener exactamente 3 caracteres"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestLen_String_FailsWhenLonger(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Len("code", "abcd", map[string]any{}, []string{"3"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "code", "len")
	want := "El campo code debe tener exactamente 3 caracteres"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestLen_String_Fails_WithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Len("code", "ab", map[string]any{}, []string{"3"}, "2", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "code", "len")
	want := "El campo code en la posición 2 debe tener exactamente 3 caracteres"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestLen_String_Fails_CustomError(t *testing.T) {
	errors := make(map[string]interface{})
	customErrors := map[string]string{"code.len": "mensaje personalizado"}
	errors = Len("code", "ab", map[string]any{}, []string{"3"}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "code", "len")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestLen_Slice_Passes(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Len("tags", []string{"a", "b"}, map[string]any{}, []string{"2"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestLen_SliceAny_Passes(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Len("items", []any{1, "x", true}, map[string]any{}, []string{"3"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestLen_Slice_Fails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Len("tags", []int{1, 2, 3}, map[string]any{}, []string{"2"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "tags", "len")
	want := "El campo tags debe tener exactamente 2 elementos"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestLen_Slice_Fails_WithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Len("tags", []int{1}, map[string]any{}, []string{"2"}, "0", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "tags", "len")
	want := "El campo tags en la posición 0 debe tener exactamente 2 elementos"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestLen_Array_Passes(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Len("pair", [2]int{1, 2}, map[string]any{}, []string{"2"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestLen_Array_Fails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Len("pair", [2]int{1, 2}, map[string]any{}, []string{"3"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "pair", "len")
	want := "El campo pair debe tener exactamente 3 elementos"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestLen_Map_Passes(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Len("meta", map[string]any{"a": 1, "b": 2}, map[string]any{}, []string{"2"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestLen_Map_Fails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Len("meta", map[string]any{"a": 1}, map[string]any{}, []string{"2"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "meta", "len")
	want := "El campo meta debe tener exactamente 2 elementos"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestLen_UncountableTypes_AreIgnored(t *testing.T) {
	values := []any{5, 3.14, true, nil}
	for _, v := range values {
		errors := make(map[string]interface{})
		errors = Len("field", v, map[string]any{}, []string{"2"}, "", errors, testAddError, map[string]string{})

		if len(errors) != 0 {
			t.Fatalf("value %v: expected no errors, got %v", v, errors)
		}
	}
}
