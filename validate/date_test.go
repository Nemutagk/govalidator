package validate

import "testing"

func TestDate_NilValueSkips(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Date("start_date", nil, map[string]any{}, []string{}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestDate_EmptyValueSkips(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Date("start_date", "", map[string]any{}, []string{}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestDate_DefaultLayoutPasses(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"start_date": "2024-01-15T10:30:00"}
	errors = Date("start_date", "2024-01-15T10:30:00", payload, []string{}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestDate_FieldMissingFromPayloadFails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Date("start_date", "2024-01-15T10:30:00", map[string]any{}, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "start_date", "date")
	want := `El campo "start_date" debe ser una fecha válida con el formato "2006-01-02T15:04:05"`
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

// payload[input] no es string aunque value sí lo sea: la validación de existencia
// se hace contra payload, no contra value, así que también falla en este caso.
func TestDate_FieldNotStringInPayloadFails(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"start_date": 20240115}
	errors = Date("start_date", 20240115, payload, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "start_date", "date")
	want := `El campo "start_date" debe ser una fecha válida con el formato "2006-01-02T15:04:05"`
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestDate_MalformedValueFails(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"start_date": "not-a-date"}
	errors = Date("start_date", "not-a-date", payload, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "start_date", "date")
	want := `El campo "start_date" debe ser una fecha válida con el formato "2006-01-02T15:04:05"`
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestDate_CustomLayoutPasses(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"birthdate": "2024-01-15"}
	errors = Date("birthdate", "2024-01-15", payload, []string{"2006-01-02"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestDate_CustomLayoutFails(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"birthdate": "15/01/2024"}
	errors = Date("birthdate", "15/01/2024", payload, []string{"2006-01-02"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "birthdate", "date")
	want := `El campo "birthdate" debe ser una fecha válida con el formato "2006-01-02"`
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestDate_EmptyLayoutOptionUsesDefault(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"start_date": "2024-01-15T10:30:00"}
	errors = Date("start_date", "2024-01-15T10:30:00", payload, []string{""}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestDate_WithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"start_date": "not-a-date"}
	errors = Date("start_date", "not-a-date", payload, []string{}, "2", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "start_date", "date")
	want := `El campo "start_date" en la posición 2 debe ser una fecha válida con el formato "2006-01-02T15:04:05"`
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestDate_MissingFieldWithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Date("items", "2024-01-15T10:30:00", map[string]any{}, []string{}, "1", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "items", "date")
	want := `El campo "items" en la posición 1 debe ser una fecha válida con el formato "2006-01-02T15:04:05"`
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestDate_CustomErrorMessage(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"start_date": "not-a-date"}
	customErrors := map[string]string{"start_date.date": "mensaje personalizado"}
	errors = Date("start_date", "not-a-date", payload, []string{}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "start_date", "date")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestDate_CustomErrorMessageOnMissingField(t *testing.T) {
	errors := make(map[string]interface{})
	customErrors := map[string]string{"start_date.date": "mensaje personalizado"}
	errors = Date("start_date", "2024-01-15T10:30:00", map[string]any{}, []string{}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "start_date", "date")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}
