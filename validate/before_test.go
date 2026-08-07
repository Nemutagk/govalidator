package validate

import (
	"testing"
	"time"
)

func TestBefore_NoOptionsFails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Before("age", "2024-01-01", map[string]any{}, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "age", "before")
	if !found || msgs[0] != "El valor a comparar no está definido" {
		t.Fatalf("msgs = %v, want 'El valor a comparar no está definido'", msgs)
	}
}

func TestBefore_NilValueSkips(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Before("end_date", nil, map[string]any{}, []string{"2024-01-01"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestBefore_EmptyValueSkips(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Before("end_date", "", map[string]any{}, []string{"2024-01-01"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestBefore_MalformedOwnValueFails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Before("end_date", "not-a-date", map[string]any{}, []string{"2024-01-01"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "end_date", "before")
	want := "La fecha proporcionada no es válida o no coincide con el formato 2006-01-02"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestBefore_IntComparisonPasses(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Before("age", 3, map[string]any{}, []string{"5"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestBefore_IntComparisonFails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Before("age", 10, map[string]any{}, []string{"5"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "age", "before")
	want := "La entrada age no es anterior al número 5"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestBefore_IntComparisonWithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Before("age", 10, map[string]any{}, []string{"5"}, "2", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "age", "before")
	want := "La entrada en la posición 2 no es anterior al número 5"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestBefore_IntComparisonCustomErrorMessage(t *testing.T) {
	errors := make(map[string]interface{})
	customErrors := map[string]string{"age.before": "mensaje personalizado"}
	errors = Before("age", 10, map[string]any{}, []string{"5"}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "age", "before")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestBefore_AgainstPayloadFieldPasses(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"end_date": "2024-02-01"}
	errors = Before("start_date", "2024-01-01", payload, []string{"end_date"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestBefore_AgainstPayloadFieldFails(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"end_date": "2024-01-01"}
	errors = Before("start_date", "2024-02-01", payload, []string{"end_date"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "start_date", "before")
	want := "La fecha no es anterior a la fecha 2024-01-01"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestBefore_AgainstPayloadFieldEqualFails(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"end_date": "2024-01-01"}
	errors = Before("start_date", "2024-01-01", payload, []string{"end_date"}, "", errors, testAddError, map[string]string{})

	if _, found := getErrorMsgs(errors, "start_date", "before"); !found {
		t.Fatalf("expected equal dates to fail, got %v", errors)
	}
}

func TestBefore_AgainstPayloadFieldInvalidDateFails(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"end_date": "not-a-date"}
	errors = Before("start_date", "2024-01-01", payload, []string{"end_date"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "start_date", "before")
	want := "La fecha de comparación no es válida o no coincide con el formato 2006-01-02"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestBefore_AgainstPayloadFieldCustomFormat(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"end_date": "01/02/2024"}
	errors = Before("start_date", "01/01/2024", payload, []string{"end_date", "02/01/2006"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestBefore_AgainstPayloadFieldCustomErrorMessage(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"end_date": "2024-01-01"}
	customErrors := map[string]string{"start_date.before": "mensaje personalizado"}
	errors = Before("start_date", "2024-02-01", payload, []string{"end_date"}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "start_date", "before")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestBefore_TodayKeywordPasses(t *testing.T) {
	errors := make(map[string]interface{})
	past := time.Now().AddDate(-1, 0, 0).Format("2006-01-02")
	errors = Before("start_date", past, map[string]any{}, []string{"today"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestBefore_TodayKeywordFails(t *testing.T) {
	errors := make(map[string]interface{})
	future := time.Now().AddDate(1, 0, 0).Format("2006-01-02")
	errors = Before("start_date", future, map[string]any{}, []string{"today"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "start_date", "before")
	want := "La fecha no es anterior a la fecha actual"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestBefore_TomorrowKeywordPasses(t *testing.T) {
	errors := make(map[string]interface{})
	past := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	errors = Before("start_date", past, map[string]any{}, []string{"tomorrow"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestBefore_TomorrowKeywordFails(t *testing.T) {
	errors := make(map[string]interface{})
	future := time.Now().AddDate(0, 0, 5).Format("2006-01-02")
	errors = Before("start_date", future, map[string]any{}, []string{"tomorrow"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "start_date", "before")
	want := "La fecha no es anterior a mañana"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestBefore_YesterdayKeywordPasses(t *testing.T) {
	errors := make(map[string]interface{})
	past := time.Now().AddDate(0, 0, -5).Format("2006-01-02")
	errors = Before("start_date", past, map[string]any{}, []string{"yesterday"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestBefore_YesterdayKeywordFails(t *testing.T) {
	errors := make(map[string]interface{})
	future := time.Now().AddDate(0, 0, 5).Format("2006-01-02")
	errors = Before("start_date", future, map[string]any{}, []string{"yesterday"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "start_date", "before")
	want := "La fecha no es anterior a ayer"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestBefore_KeywordCustomErrorMessage(t *testing.T) {
	errors := make(map[string]interface{})
	future := time.Now().AddDate(1, 0, 0).Format("2006-01-02")
	customErrors := map[string]string{"start_date.before": "mensaje personalizado"}
	errors = Before("start_date", future, map[string]any{}, []string{"today"}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "start_date", "before")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestBefore_LiteralDatePasses(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Before("start_date", "2024-01-01", map[string]any{}, []string{"2024-02-01"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestBefore_LiteralDateFails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Before("start_date", "2024-02-01", map[string]any{}, []string{"2024-01-01"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "start_date", "before")
	want := "La fecha no es anterior a la fecha 2024-01-01"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestBefore_LiteralDateCustomErrorMessage(t *testing.T) {
	errors := make(map[string]interface{})
	customErrors := map[string]string{"start_date.before": "mensaje personalizado"}
	errors = Before("start_date", "2024-02-01", map[string]any{}, []string{"2024-01-01"}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "start_date", "before")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

// A diferencia de After (que ignora en silencio una opción literal no parseable),
// Before sí agrega un error explícito cuando options[0] no es un campo del
// payload, ni una palabra clave humana, ni una fecha literal válida.
func TestBefore_UnparseableLiteralOptionFails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Before("start_date", "2024-01-01", map[string]any{}, []string{"not-a-real-option"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "start_date", "before")
	want := "La fecha de comparación no es válida o no coincide con el formato 2006-01-02"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}
