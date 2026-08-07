package validate

import (
	"testing"
	"time"
)

func TestToFloat64(t *testing.T) {
	cases := []struct {
		name    string
		value   any
		want    float64
		wantOk  bool
	}{
		{"int", 10, 10, true},
		{"int64", int64(20), 20, true},
		{"float64", 3.5, 3.5, true},
		{"string not supported", "10", 0, false},
		{"bool not supported", true, 0, false},
		{"nil not supported", nil, 0, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, ok := toFloat64(c.value)
			if ok != c.wantOk {
				t.Fatalf("ok = %v, want %v", ok, c.wantOk)
			}
			if ok && got != c.want {
				t.Fatalf("got = %v, want %v", got, c.want)
			}
		})
	}
}

func TestParseFloat(t *testing.T) {
	if v, err := parseFloat("3.14"); err != nil || v != 3.14 {
		t.Fatalf("parseFloat(3.14) = (%v, %v), want (3.14, nil)", v, err)
	}

	if _, err := parseFloat("not-a-number"); err == nil {
		t.Fatalf("expected error parsing invalid float")
	}
}

func TestResolveComparable_NumericLiteral(t *testing.T) {
	errors := make(map[string]interface{})
	pair, errors, ok := resolveComparable("age", 10, map[string]any{}, []string{"5"}, "greater_than", "", errors, testAddError, map[string]string{})

	if !ok {
		t.Fatalf("expected ok=true, errors=%v", errors)
	}
	if pair.isDate {
		t.Fatalf("expected numeric comparison, got isDate=true")
	}
	if pair.ownNum != 10 || pair.cmpNum != 5 {
		t.Fatalf("ownNum=%v cmpNum=%v, want 10 and 5", pair.ownNum, pair.cmpNum)
	}
}

func TestResolveComparable_NumericAgainstField(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"min": 7}
	pair, errors, ok := resolveComparable("age", 10, payload, []string{"min"}, "greater_than", "", errors, testAddError, map[string]string{})

	if !ok {
		t.Fatalf("expected ok=true, errors=%v", errors)
	}
	if pair.ownNum != 10 || pair.cmpNum != 7 {
		t.Fatalf("ownNum=%v cmpNum=%v, want 10 and 7", pair.ownNum, pair.cmpNum)
	}
}

func TestResolveComparable_NumericAgainstFieldInvalid(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"min": "not-a-number"}
	_, errors, ok := resolveComparable("age", 10, payload, []string{"min"}, "greater_than", "", errors, testAddError, map[string]string{})

	if ok {
		t.Fatalf("expected ok=false")
	}
	msgs, found := getErrorMsgs(errors, "age", "greater_than")
	if !found || len(msgs) != 1 {
		t.Fatalf("expected one error message, got %v", errors)
	}
}

func TestResolveComparable_NumericLiteralInvalid(t *testing.T) {
	errors := make(map[string]interface{})
	_, errors, ok := resolveComparable("age", 10, map[string]any{}, []string{"not-a-number"}, "greater_than", "", errors, testAddError, map[string]string{})

	if ok {
		t.Fatalf("expected ok=false")
	}
	if _, found := getErrorMsgs(errors, "age", "greater_than"); !found {
		t.Fatalf("expected error message, got %v", errors)
	}
}

func TestResolveComparable_ValueNotNumberOrString(t *testing.T) {
	errors := make(map[string]interface{})
	_, errors, ok := resolveComparable("age", true, map[string]any{}, []string{"5"}, "greater_than", "", errors, testAddError, map[string]string{})

	if ok {
		t.Fatalf("expected ok=false")
	}
	if _, found := getErrorMsgs(errors, "age", "greater_than"); !found {
		t.Fatalf("expected error message, got %v", errors)
	}
}

func TestResolveComparable_DateLiteral(t *testing.T) {
	errors := make(map[string]interface{})
	pair, errors, ok := resolveComparable("birthdate", "2024-01-02", map[string]any{}, []string{"2024-01-01"}, "greater_than", "", errors, testAddError, map[string]string{})

	if !ok {
		t.Fatalf("expected ok=true, errors=%v", errors)
	}
	if !pair.isDate {
		t.Fatalf("expected date comparison, got isDate=false")
	}

	wantOwn, _ := time.Parse(defaultCompareDateLayout, "2024-01-02")
	wantCmp, _ := time.Parse(defaultCompareDateLayout, "2024-01-01")
	if !pair.ownDate.Equal(wantOwn) || !pair.cmpDate.Equal(wantCmp) {
		t.Fatalf("ownDate=%v cmpDate=%v, want %v and %v", pair.ownDate, pair.cmpDate, wantOwn, wantCmp)
	}
}

func TestResolveComparable_DateAgainstField(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"start_date": "2024-01-01"}
	pair, errors, ok := resolveComparable("end_date", "2024-01-02", payload, []string{"start_date"}, "greater_than", "", errors, testAddError, map[string]string{})

	if !ok {
		t.Fatalf("expected ok=true, errors=%v", errors)
	}
	if !pair.isDate {
		t.Fatalf("expected date comparison")
	}
}

func TestResolveComparable_DateAgainstFieldNotString(t *testing.T) {
	errors := make(map[string]interface{})
	payload := map[string]any{"start_date": 123}
	_, errors, ok := resolveComparable("end_date", "2024-01-02", payload, []string{"start_date"}, "greater_than", "", errors, testAddError, map[string]string{})

	if ok {
		t.Fatalf("expected ok=false")
	}
	if _, found := getErrorMsgs(errors, "end_date", "greater_than"); !found {
		t.Fatalf("expected error message, got %v", errors)
	}
}

func TestResolveComparable_OwnDateInvalid(t *testing.T) {
	errors := make(map[string]interface{})
	_, errors, ok := resolveComparable("birthdate", "not-a-date", map[string]any{}, []string{"2024-01-01"}, "greater_than", "", errors, testAddError, map[string]string{})

	if ok {
		t.Fatalf("expected ok=false")
	}
	if _, found := getErrorMsgs(errors, "birthdate", "greater_than"); !found {
		t.Fatalf("expected error message, got %v", errors)
	}
}

func TestResolveComparable_OwnDateInvalidWithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	_, errors, ok := resolveComparable("birthdate", "not-a-date", map[string]any{}, []string{"2024-01-01"}, "greater_than", "2", errors, testAddError, map[string]string{})

	if ok {
		t.Fatalf("expected ok=false")
	}
	msgs, found := getErrorMsgs(errors, "birthdate", "greater_than")
	if !found || len(msgs) != 1 {
		t.Fatalf("expected error message, got %v", errors)
	}
	want := "El campo birthdate en la posición 2 debe ser un número o una fecha válida con el formato \"2006-01-02\""
	if msgs[0] != want {
		t.Fatalf("msg = %q, want %q", msgs[0], want)
	}
}

func TestResolveComparable_TargetDateInvalid(t *testing.T) {
	errors := make(map[string]interface{})
	_, errors, ok := resolveComparable("birthdate", "2024-01-01", map[string]any{}, []string{"not-a-date"}, "greater_than", "", errors, testAddError, map[string]string{})

	if ok {
		t.Fatalf("expected ok=false")
	}
	if _, found := getErrorMsgs(errors, "birthdate", "greater_than"); !found {
		t.Fatalf("expected error message, got %v", errors)
	}
}

func TestResolveComparable_CustomLayouts(t *testing.T) {
	errors := make(map[string]interface{})
	// options[1] = layout propio, options[2] = layout del valor comparado
	pair, errors, ok := resolveComparable("birthdate", "03/01/2024", map[string]any{}, []string{"2024-01-02", "02/01/2006", "2006-01-02"}, "greater_than", "", errors, testAddError, map[string]string{})

	if !ok {
		t.Fatalf("expected ok=true, errors=%v", errors)
	}

	wantOwn, _ := time.Parse("02/01/2006", "03/01/2024")
	wantCmp, _ := time.Parse("2006-01-02", "2024-01-02")
	if !pair.ownDate.Equal(wantOwn) {
		t.Fatalf("ownDate=%v want %v", pair.ownDate, wantOwn)
	}
	if !pair.cmpDate.Equal(wantCmp) {
		t.Fatalf("cmpDate=%v want %v", pair.cmpDate, wantCmp)
	}
}

func TestResolveComparable_CustomErrorMessage(t *testing.T) {
	errors := make(map[string]interface{})
	customErrors := map[string]string{"age.greater_than": "mensaje personalizado"}
	_, errors, ok := resolveComparable("age", true, map[string]any{}, []string{"5"}, "greater_than", "", errors, testAddError, customErrors)

	if ok {
		t.Fatalf("expected ok=false")
	}
	msgs, found := getErrorMsgs(errors, "age", "greater_than")
	if !found || len(msgs) != 1 || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}
