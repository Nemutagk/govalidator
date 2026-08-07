package validate

import "testing"

func TestGreaterThanEqual_NoOptions(t *testing.T) {
	errors := make(map[string]interface{})
	errors = GreaterThanEqual("age", 10, map[string]any{}, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "age", "greater_than_equal")
	if !found || msgs[0] != "El valor a comparar no está definido" {
		t.Fatalf("expected 'no definido' error, got %v", errors)
	}
}

func TestGreaterThanEqual_NumericGreaterPasses(t *testing.T) {
	errors := make(map[string]interface{})
	errors = GreaterThanEqual("age", 10, map[string]any{}, []string{"5"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestGreaterThanEqual_NumericEqualPasses(t *testing.T) {
	errors := make(map[string]interface{})
	errors = GreaterThanEqual("age", 5, map[string]any{}, []string{"5"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected equal values to pass, got %v", errors)
	}
}

func TestGreaterThanEqual_NumericLessFails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = GreaterThanEqual("age", 3, map[string]any{}, []string{"5"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "age", "greater_than_equal")
	if !found || msgs[0] != "El campo age debe ser mayor o igual que 5" {
		t.Fatalf("expected lesser value to fail, got %v", errors)
	}
}

func TestGreaterThanEqual_WithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	errors = GreaterThanEqual("age", 3, map[string]any{}, []string{"5"}, "2", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "age", "greater_than_equal")
	want := "El campo age en la posición 2 debe ser mayor o igual que 5"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestGreaterThanEqual_CustomErrorMessage(t *testing.T) {
	errors := make(map[string]interface{})
	customErrors := map[string]string{"age.greater_than_equal": "mensaje personalizado"}
	errors = GreaterThanEqual("age", 3, map[string]any{}, []string{"5"}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "age", "greater_than_equal")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestGreaterThanEqual_DateEqualPasses(t *testing.T) {
	errors := make(map[string]interface{})
	errors = GreaterThanEqual("birthdate", "2024-01-01", map[string]any{}, []string{"2024-01-01"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected equal dates to pass, got %v", errors)
	}
}

func TestGreaterThanEqual_DateEarlierFails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = GreaterThanEqual("birthdate", "2023-12-31", map[string]any{}, []string{"2024-01-01"}, "", errors, testAddError, map[string]string{})

	if _, found := getErrorMsgs(errors, "birthdate", "greater_than_equal"); !found {
		t.Fatalf("expected earlier date to fail, got %v", errors)
	}
}

func TestGreaterThanEqual_DelegatesResolveFailure(t *testing.T) {
	errors := make(map[string]interface{})
	errors = GreaterThanEqual("age", true, map[string]any{}, []string{"5"}, "", errors, testAddError, map[string]string{})

	if _, found := getErrorMsgs(errors, "age", "greater_than_equal"); !found {
		t.Fatalf("expected resolveComparable failure to propagate, got %v", errors)
	}
}
