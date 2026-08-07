package validate

import "testing"

func TestGreaterThan_NoOptions(t *testing.T) {
	errors := make(map[string]interface{})
	errors = GreaterThan("age", 10, map[string]any{}, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "age", "greater_than")
	if !found || len(msgs) != 1 || msgs[0] != "El valor a comparar no está definido" {
		t.Fatalf("expected 'no definido' error, got %v", errors)
	}
}

func TestGreaterThan_NumericPasses(t *testing.T) {
	errors := make(map[string]interface{})
	errors = GreaterThan("age", 10, map[string]any{}, []string{"5"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestGreaterThan_NumericEqualFails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = GreaterThan("age", 5, map[string]any{}, []string{"5"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "age", "greater_than")
	if !found || msgs[0] != "El campo age debe ser mayor que 5" {
		t.Fatalf("expected equal value to fail, got %v", errors)
	}
}

func TestGreaterThan_NumericLessFails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = GreaterThan("age", 3, map[string]any{}, []string{"5"}, "", errors, testAddError, map[string]string{})

	if _, found := getErrorMsgs(errors, "age", "greater_than"); !found {
		t.Fatalf("expected lesser value to fail, got %v", errors)
	}
}

func TestGreaterThan_WithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	errors = GreaterThan("age", 3, map[string]any{}, []string{"5"}, "2", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "age", "greater_than")
	want := "El campo age en la posición 2 debe ser mayor que 5"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestGreaterThan_CustomErrorMessage(t *testing.T) {
	errors := make(map[string]interface{})
	customErrors := map[string]string{"age.greater_than": "mensaje personalizado"}
	errors = GreaterThan("age", 3, map[string]any{}, []string{"5"}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "age", "greater_than")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestGreaterThan_DatePasses(t *testing.T) {
	errors := make(map[string]interface{})
	errors = GreaterThan("birthdate", "2024-01-02", map[string]any{}, []string{"2024-01-01"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestGreaterThan_DateEqualFails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = GreaterThan("birthdate", "2024-01-01", map[string]any{}, []string{"2024-01-01"}, "", errors, testAddError, map[string]string{})

	if _, found := getErrorMsgs(errors, "birthdate", "greater_than"); !found {
		t.Fatalf("expected equal dates to fail, got %v", errors)
	}
}

func TestGreaterThan_DateEarlierFails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = GreaterThan("birthdate", "2023-12-31", map[string]any{}, []string{"2024-01-01"}, "", errors, testAddError, map[string]string{})

	if _, found := getErrorMsgs(errors, "birthdate", "greater_than"); !found {
		t.Fatalf("expected earlier date to fail, got %v", errors)
	}
}

func TestGreaterThan_DelegatesResolveFailure(t *testing.T) {
	errors := make(map[string]interface{})
	errors = GreaterThan("age", true, map[string]any{}, []string{"5"}, "", errors, testAddError, map[string]string{})

	if _, found := getErrorMsgs(errors, "age", "greater_than"); !found {
		t.Fatalf("expected resolveComparable failure to propagate, got %v", errors)
	}
}
