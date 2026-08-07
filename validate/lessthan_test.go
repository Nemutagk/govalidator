package validate

import "testing"

func TestLessThan_NoOptions(t *testing.T) {
	errors := make(map[string]interface{})
	errors = LessThan("age", 3, map[string]any{}, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "age", "less_than")
	if !found || msgs[0] != "El valor a comparar no está definido" {
		t.Fatalf("expected 'no definido' error, got %v", errors)
	}
}

func TestLessThan_NumericPasses(t *testing.T) {
	errors := make(map[string]interface{})
	errors = LessThan("age", 3, map[string]any{}, []string{"5"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestLessThan_NumericEqualFails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = LessThan("age", 5, map[string]any{}, []string{"5"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "age", "less_than")
	if !found || msgs[0] != "El campo age debe ser menor que 5" {
		t.Fatalf("expected equal value to fail, got %v", errors)
	}
}

func TestLessThan_NumericGreaterFails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = LessThan("age", 10, map[string]any{}, []string{"5"}, "", errors, testAddError, map[string]string{})

	if _, found := getErrorMsgs(errors, "age", "less_than"); !found {
		t.Fatalf("expected greater value to fail, got %v", errors)
	}
}

func TestLessThan_WithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	errors = LessThan("age", 10, map[string]any{}, []string{"5"}, "2", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "age", "less_than")
	want := "El campo age en la posición 2 debe ser menor que 5"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestLessThan_CustomErrorMessage(t *testing.T) {
	errors := make(map[string]interface{})
	customErrors := map[string]string{"age.less_than": "mensaje personalizado"}
	errors = LessThan("age", 10, map[string]any{}, []string{"5"}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "age", "less_than")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestLessThan_DatePasses(t *testing.T) {
	errors := make(map[string]interface{})
	errors = LessThan("birthdate", "2023-12-31", map[string]any{}, []string{"2024-01-01"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestLessThan_DateEqualFails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = LessThan("birthdate", "2024-01-01", map[string]any{}, []string{"2024-01-01"}, "", errors, testAddError, map[string]string{})

	if _, found := getErrorMsgs(errors, "birthdate", "less_than"); !found {
		t.Fatalf("expected equal dates to fail, got %v", errors)
	}
}

func TestLessThan_DateLaterFails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = LessThan("birthdate", "2024-01-02", map[string]any{}, []string{"2024-01-01"}, "", errors, testAddError, map[string]string{})

	if _, found := getErrorMsgs(errors, "birthdate", "less_than"); !found {
		t.Fatalf("expected later date to fail, got %v", errors)
	}
}

func TestLessThan_DelegatesResolveFailure(t *testing.T) {
	errors := make(map[string]interface{})
	errors = LessThan("age", true, map[string]any{}, []string{"5"}, "", errors, testAddError, map[string]string{})

	if _, found := getErrorMsgs(errors, "age", "less_than"); !found {
		t.Fatalf("expected resolveComparable failure to propagate, got %v", errors)
	}
}
