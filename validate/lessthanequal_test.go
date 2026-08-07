package validate

import "testing"

func TestLessThanEqual_NoOptions(t *testing.T) {
	errors := make(map[string]interface{})
	errors = LessThanEqual("age", 3, map[string]any{}, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "age", "less_than_equal")
	if !found || msgs[0] != "El valor a comparar no está definido" {
		t.Fatalf("expected 'no definido' error, got %v", errors)
	}
}

func TestLessThanEqual_NumericLessPasses(t *testing.T) {
	errors := make(map[string]interface{})
	errors = LessThanEqual("age", 3, map[string]any{}, []string{"5"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestLessThanEqual_NumericEqualPasses(t *testing.T) {
	errors := make(map[string]interface{})
	errors = LessThanEqual("age", 5, map[string]any{}, []string{"5"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected equal values to pass, got %v", errors)
	}
}

func TestLessThanEqual_NumericGreaterFails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = LessThanEqual("age", 10, map[string]any{}, []string{"5"}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "age", "less_than_equal")
	if !found || msgs[0] != "El campo age debe ser menor o igual que 5" {
		t.Fatalf("expected greater value to fail, got %v", errors)
	}
}

func TestLessThanEqual_WithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	errors = LessThanEqual("age", 10, map[string]any{}, []string{"5"}, "2", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "age", "less_than_equal")
	want := "El campo age en la posición 2 debe ser menor o igual que 5"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestLessThanEqual_CustomErrorMessage(t *testing.T) {
	errors := make(map[string]interface{})
	customErrors := map[string]string{"age.less_than_equal": "mensaje personalizado"}
	errors = LessThanEqual("age", 10, map[string]any{}, []string{"5"}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "age", "less_than_equal")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestLessThanEqual_DateEqualPasses(t *testing.T) {
	errors := make(map[string]interface{})
	errors = LessThanEqual("birthdate", "2024-01-01", map[string]any{}, []string{"2024-01-01"}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected equal dates to pass, got %v", errors)
	}
}

func TestLessThanEqual_DateLaterFails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = LessThanEqual("birthdate", "2024-01-02", map[string]any{}, []string{"2024-01-01"}, "", errors, testAddError, map[string]string{})

	if _, found := getErrorMsgs(errors, "birthdate", "less_than_equal"); !found {
		t.Fatalf("expected later date to fail, got %v", errors)
	}
}

func TestLessThanEqual_DelegatesResolveFailure(t *testing.T) {
	errors := make(map[string]interface{})
	errors = LessThanEqual("age", true, map[string]any{}, []string{"5"}, "", errors, testAddError, map[string]string{})

	if _, found := getErrorMsgs(errors, "age", "less_than_equal"); !found {
		t.Fatalf("expected resolveComparable failure to propagate, got %v", errors)
	}
}
