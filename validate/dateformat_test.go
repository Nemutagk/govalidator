package validate

import "testing"

func TestDateFormat_NoOptionsFails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = DateFormat("start_date", "2024-01-15", map[string]any{}, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "start_date", "date_format")
	if !found || msgs[0] != "El formato de fecha no está definido" {
		t.Fatalf("msgs = %v, want 'El formato de fecha no está definido'", msgs)
	}
}

// La rama de "no options" no consulta customeErrors, siempre usa el mensaje fijo.
func TestDateFormat_NoOptionsIgnoresCustomError(t *testing.T) {
	errors := make(map[string]interface{})
	customErrors := map[string]string{"start_date.date_format": "mensaje personalizado"}
	errors = DateFormat("start_date", "2024-01-15", map[string]any{}, []string{}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "start_date", "date_format")
	if !found || msgs[0] != "El formato de fecha no está definido" {
		t.Fatalf("expected default message even with customeErrors set, got %v", msgs)
	}
}

func TestDateFormat_NilValueSkips(t *testing.T) {
	errors := make(map[string]interface{})
	errors = DateFormat("start_date", nil, map[string]any{}, []string{"2006-01-02"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestDateFormat_EmptyValueSkips(t *testing.T) {
	errors := make(map[string]interface{})
	errors = DateFormat("start_date", "", map[string]any{}, []string{"2006-01-02"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestDateFormat_ValidFormatPasses(t *testing.T) {
	errors := make(map[string]interface{})
	errors = DateFormat("start_date", "2024-01-15", map[string]any{}, []string{"2006-01-02"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestDateFormat_InvalidFormatFails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = DateFormat("start_date", "15/01/2024", map[string]any{}, []string{"2006-01-02"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "start_date", "date_format")
	if !found || msgs[0] != "El formato de la fecha es inválido" {
		t.Fatalf("msgs = %v, want 'El formato de la fecha es inválido'", msgs)
	}
}

func TestDateFormat_CustomFormatLayout(t *testing.T) {
	errors := make(map[string]interface{})
	errors = DateFormat("start_date", "15/01/2024", map[string]any{}, []string{"02/01/2006"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestDateFormat_WithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	errors = DateFormat("items", "not-a-date", map[string]any{}, []string{"2006-01-02"}, "3", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "items", "date_format")
	want := "El formato de la fecha en el índice 3 es inválido"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestDateFormat_CustomErrorMessage(t *testing.T) {
	errors := make(map[string]interface{})
	customErrors := map[string]string{"start_date.date_format": "mensaje personalizado"}
	errors = DateFormat("start_date", "not-a-date", map[string]any{}, []string{"2006-01-02"}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "start_date", "date_format")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestDateFormat_NonStringValueFails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = DateFormat("age", 20240115, map[string]any{}, []string{"2006-01-02"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "age", "date_format")
	want := "El campo \"age\" debe ser una fecha válida con el formato \"2006-01-02\""
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestDateFormat_NonStringValueWithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	errors = DateFormat("age", 20240115, map[string]any{}, []string{"2006-01-02"}, "3", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "age", "date_format")
	want := "El campo \"age\" en la posición 3 debe ser una fecha válida con el formato \"2006-01-02\""
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestDateFormat_NonStringValueCustomErrorMessage(t *testing.T) {
	errors := make(map[string]interface{})
	customErrors := map[string]string{"age.date_format": "mensaje personalizado"}
	errors = DateFormat("age", 20240115, map[string]any{}, []string{"2006-01-02"}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "age", "date_format")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}
