package govalidator

import (
	"strings"
	"testing"
)

type typedMaps struct {
	Name   map[string]string `json:"name"`
	Scores map[int]int       `json:"scores"`
	Labels map[int]string    `json:"labels"`
	Flags  map[uint8]bool    `json:"flags"`
}

func required(name string) Input {
	return Input{Name: name, Rules: []Rule{{Name: "required"}}}
}

func errorsOf(t *testing.T, err error) []string {
	t.Helper()
	if err == nil {
		return nil
	}
	ve, ok := AsValidationError(err)
	if !ok {
		t.Fatalf("expected a validation error, got %v", err)
	}
	return ve.Errors
}

func TestValidateStruct_TypedMaps_DottedRules(t *testing.T) {
	req := typedMaps{
		Name:   map[string]string{"es": "hola"},
		Scores: map[int]int{1: 15, 2: 7},
		Labels: map[int]string{2: "a"},
		Flags:  map[uint8]bool{3: true},
	}

	rules := []Input{
		required("name.es"),
		{Name: "scores.1", Rules: []Rule{{Name: "min", Options: []string{"10"}}}},
		{Name: "labels.2", Rules: []Rule{{Name: "in", Options: []string{"a", "b"}}}},
		required("flags.3"),
	}

	result, err := ValidateStruct(req, rules, nil, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Name["es"] != "hola" || result.Scores[1] != 15 || result.Labels[2] != "a" || !result.Flags[3] {
		t.Fatalf("validated values were not preserved: %+v", result)
	}
	if _, kept := result.Scores[2]; kept {
		t.Fatalf("scores.2 has no rule and must be dropped: %+v", result.Scores)
	}
}

func TestValidateStruct_TypedMaps_Rejections(t *testing.T) {
	cases := []struct {
		name  string
		req   typedMaps
		rule  Input
		field string
	}{
		{"missing string key", typedMaps{Name: map[string]string{"en": "x"}}, required("name.es"), "name.es"},
		{"empty map", typedMaps{Name: map[string]string{}}, required("name.es"), "name.es"},
		{"nil map", typedMaps{}, required("name.es"), "name.es"},
		{"missing int key", typedMaps{Scores: map[int]int{2: 1}}, required("scores.1"), "scores.1"},
		{"int value below min", typedMaps{Scores: map[int]int{1: 5}}, Input{Name: "scores.1", Rules: []Rule{{Name: "min", Options: []string{"10"}}}}, "scores.1"},
		{"string outside in", typedMaps{Labels: map[int]string{2: "z"}}, Input{Name: "labels.2", Rules: []Rule{{Name: "in", Options: []string{"a", "b"}}}}, "labels.2"},
		{"empty string value", typedMaps{Name: map[string]string{"es": ""}}, required("name.es"), "name.es"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := ValidateStruct(c.req, []Input{c.rule}, nil, nil)
			errs := errorsOf(t, err)
			if len(errs) == 0 {
				t.Fatalf("expected a validation error")
			}
			if !strings.Contains(errs[0], c.field) {
				t.Fatalf("error %q does not mention %q", errs[0], c.field)
			}
		})
	}
}

func TestValidateStruct_TypedMaps_BareRulesSeeTheOriginalMap(t *testing.T) {
	req := typedMaps{Name: map[string]string{"es": "hola"}, Scores: map[int]int{1: 1}}

	rules := []Input{
		{Name: "name", Rules: []Rule{{Name: "required"}, {Name: "type", Options: []string{"map[string]string"}}}},
		{Name: "scores", Rules: []Rule{{Name: "type", Options: []string{"map[int]int"}}, {Name: "len", Options: []string{"1"}}}},
	}

	if _, err := ValidateStruct(req, rules, nil, nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateStruct_SliceOfTypedMaps(t *testing.T) {
	type payload struct {
		Items []map[string]string `json:"items"`
	}

	rules := []Input{required("items.*.name")}

	if _, err := ValidateStruct(payload{Items: []map[string]string{{"name": "a"}, {"name": "b"}}}, rules, nil, nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err := ValidateStruct(payload{Items: []map[string]string{{"name": "a"}, {"other": "b"}}}, rules, nil, nil)
	if len(errorsOf(t, err)) == 0 {
		t.Fatalf("expected an error for the item without name")
	}
}

func TestValidateStruct_RequiredIfPathThroughTypedMap(t *testing.T) {
	type payload struct {
		Meta   map[string]string `json:"meta"`
		Detail string            `json:"detail"`
	}

	rules := []Input{
		{Name: "detail", Rules: []Rule{{Name: "sometimes"}, {Name: "required_if", Options: []string{"meta.mode", "batch"}}}},
		required("meta.mode"),
	}

	_, err := ValidateStruct(payload{Meta: map[string]string{"mode": "batch"}}, rules, nil, nil)
	if len(errorsOf(t, err)) == 0 {
		t.Fatalf("expected detail to be required when meta.mode is batch")
	}

	if _, err := ValidateStruct(payload{Meta: map[string]string{"mode": "batch"}, Detail: "x"}, rules, nil, nil); err != nil {
		t.Fatalf("expected no error with detail, got %v", err)
	}

	if _, err := ValidateStruct(payload{Meta: map[string]string{"mode": "single"}}, rules, nil, nil); err != nil {
		t.Fatalf("expected no error when meta.mode is not batch, got %v", err)
	}
}

func TestValidateRequest_NormalizersStillWriteBackIntoMapStringAny(t *testing.T) {
	// map[string]any no se copia: el normalizer escribe en el elemento original del body.
	elem := map[string]any{"mode": "  A  "}
	body := map[string]any{"items": []any{elem}}

	rules := []Input{{
		Name:        "items.*.mode",
		Normalizers: []Normalizer{{Name: "trim"}, {Name: "lower"}},
		Rules:       []Rule{{Name: "required"}},
	}}

	if _, err := ValidateRequest(body, rules, nil, nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if elem["mode"] != "a" {
		t.Fatalf("normalized value was not written back to the original element: %v", elem["mode"])
	}
}

func TestValidateRequest_MapWithUnsupportedKeyTypeIsNotRecorrido(t *testing.T) {
	body := map[string]any{"m": map[float64]string{1.5: "x"}}

	_, err := ValidateRequest(body, []Input{required("m.1.5")}, nil, nil)
	if len(errorsOf(t, err)) == 0 {
		t.Fatalf("a map with float keys must not be addressable by path")
	}
}
