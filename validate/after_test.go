package validate

import (
	"testing"
	"time"
)

func TestAfter_NoOptionsFails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = After("age", "2024-01-01", map[string]any{}, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "age", "after")
	if !found || msgs[0] != "El valor a comprar no esta definido" {
		t.Fatalf("msgs = %v, want 'El valor a comprar no esta definido'", msgs)
	}
}

func TestAfter_NilValueSkips(t *testing.T) {
	errors := make(map[string]interface{})
	errors = After("end_date", nil, map[string]any{}, []string{"2024-01-01"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestAfter_EmptyValueSkips(t *testing.T) {
	errors := make(map[string]interface{})
	errors = After("end_date", "", map[string]any{}, []string{"2024-01-01"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestAfter_MalformedOwnValueFails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = After("end_date", "not-a-date", map[string]any{}, []string{"2024-01-01"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "end_date", "after")
	if !found || msgs[0] != "El valor no es una fecha válida" {
		t.Fatalf("msgs = %v, want 'El valor no es una fecha válida'", msgs)
	}
}

func TestAfter_AgainstPayloadFieldPasses(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"start_date": "2024-01-01"}
	errors = After("end_date", "2024-02-01", payload, []string{"start_date"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestAfter_AgainstPayloadFieldFails(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"start_date": "2024-02-01"}
	errors = After("end_date", "2024-01-01", payload, []string{"start_date"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "end_date", "after")
	want := "La fecha no es posterior a la fecha 2024-02-01"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestAfter_AgainstPayloadFieldEqualFails(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"start_date": "2024-01-01"}
	errors = After("end_date", "2024-01-01", payload, []string{"start_date"}, "", errors, testAddError, map[string]string{})

	if _, found := getErrorMsgs(errors, "end_date", "after"); !found {
		t.Fatalf("expected equal dates to fail, got %v", errors)
	}
}

func TestAfter_AgainstPayloadFieldInvalidDateFails(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"start_date": "not-a-date"}
	errors = After("end_date", "2024-01-01", payload, []string{"start_date"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "end_date", "after")
	want := "La fecha de comparación no es válida o no coincide con el formato 2006-01-02"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestAfter_AgainstPayloadFieldCustomFormat(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"start_date": "01/01/2024"}
	errors = After("end_date", "01/02/2024", payload, []string{"start_date", "02/01/2006"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestAfter_AgainstPayloadFieldCustomErrorMessage(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"start_date": "2024-02-01"}
	customErrors := map[string]string{"end_date.after": "mensaje personalizado"}
	errors = After("end_date", "2024-01-01", payload, []string{"start_date"}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "end_date", "after")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestAfter_TodayKeywordPasses(t *testing.T) {
	errors := make(map[string]interface{})
	future := time.Now().AddDate(1, 0, 0).Format("2006-01-02")
	errors = After("end_date", future, map[string]any{}, []string{"today"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestAfter_TodayKeywordFails(t *testing.T) {
	errors := make(map[string]interface{})
	past := time.Now().AddDate(-1, 0, 0).Format("2006-01-02")
	errors = After("end_date", past, map[string]any{}, []string{"today"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "end_date", "after")
	want := "La fecha no es posterior a la fecha actual"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestAfter_TomorrowKeywordPasses(t *testing.T) {
	errors := make(map[string]interface{})
	future := time.Now().AddDate(0, 0, 5).Format("2006-01-02")
	errors = After("end_date", future, map[string]any{}, []string{"tomorrow"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestAfter_TomorrowKeywordFails(t *testing.T) {
	errors := make(map[string]interface{})
	past := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	errors = After("end_date", past, map[string]any{}, []string{"tomorrow"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "end_date", "after")
	want := "La fecha no es posterior a mañana"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestAfter_YesterdayKeywordPasses(t *testing.T) {
	errors := make(map[string]interface{})
	future := time.Now().AddDate(0, 0, 5).Format("2006-01-02")
	errors = After("end_date", future, map[string]any{}, []string{"yesterday"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestAfter_YesterdayKeywordFails(t *testing.T) {
	errors := make(map[string]interface{})
	past := time.Now().AddDate(0, 0, -5).Format("2006-01-02")
	errors = After("end_date", past, map[string]any{}, []string{"yesterday"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "end_date", "after")
	want := "La fecha no es posterior a ayer"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestAfter_KeywordCustomErrorMessage(t *testing.T) {
	errors := make(map[string]interface{})
	past := time.Now().AddDate(-1, 0, 0).Format("2006-01-02")
	customErrors := map[string]string{"end_date.after": "mensaje personalizado"}
	errors = After("end_date", past, map[string]any{}, []string{"today"}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "end_date", "after")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestAfter_LiteralDatePasses(t *testing.T) {
	errors := make(map[string]interface{})
	errors = After("end_date", "2024-02-01", map[string]any{}, []string{"2024-01-01"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestAfter_LiteralDateFails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = After("end_date", "2024-01-01", map[string]any{}, []string{"2024-02-01"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "end_date", "after")
	want := "La fecha no es posterior a la fecha 2024-02-01"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestAfter_LiteralDateCustomErrorMessage(t *testing.T) {
	errors := make(map[string]interface{})
	customErrors := map[string]string{"end_date.after": "mensaje personalizado"}
	errors = After("end_date", "2024-01-01", map[string]any{}, []string{"2024-02-01"}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "end_date", "after")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestAfter_UnparseableLiteralOptionFails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = After("end_date", "2024-01-01", map[string]any{}, []string{"not-a-real-option"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "end_date", "after")
	want := "La fecha de comparación no es válida o no coincide con el formato 2006-01-02"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestAfter_IntComparisonPasses(t *testing.T) {
	errors := make(map[string]interface{})
	errors = After("age", 10, map[string]any{}, []string{"5"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestAfter_IntComparisonFails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = After("age", 3, map[string]any{}, []string{"5"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "age", "after")
	want := "La entrada age no es posterior al número 5"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestAfter_IntComparisonWithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	errors = After("age", 3, map[string]any{}, []string{"5"}, "2", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "age", "after")
	want := "La entrada en la posición 2 no es posterior al número 5"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestAfter_IntComparisonCustomErrorMessage(t *testing.T) {
	errors := make(map[string]interface{})
	customErrors := map[string]string{"age.after": "mensaje personalizado"}
	errors = After("age", 3, map[string]any{}, []string{"5"}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "age", "after")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}
