package govalidator

import "testing"

type nestedModeSecret struct {
	Mode   string `json:"mode"`
	Secret string `json:"secret"`
}

type nestedPayload struct {
	Payload nestedModeSecret `json:"payload"`
}

type nestedPayloadPtr struct {
	Payload *nestedModeSecret `json:"payload,omitempty"`
}

type itemsPayload struct {
	Items []nestedModeSecret `json:"items"`
}

func TestValidateStruct_RequiredOnNestedStruct_PassesWhenPopulated(t *testing.T) {
	p := nestedPayload{Payload: nestedModeSecret{Mode: "individual"}}

	rules := []Input{
		{Name: "payload", Rules: []Rule{{Name: "required"}}},
	}

	if _, err := ValidateStruct(p, rules, nil, nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateStruct_RequiredOnNestedStruct_FailsWhenAbsent(t *testing.T) {
	p := nestedPayloadPtr{}

	rules := []Input{
		{Name: "payload", Rules: []Rule{{Name: "required"}}},
	}

	_, err := ValidateStruct(p, rules, nil, nil)
	ve, ok := AsValidationError(err)
	if !ok {
		t.Fatalf("expected a ValidationError, got %v", err)
	}
	want := "payload: El campo payload no está definido"
	if len(ve.Errors) != 1 || ve.Errors[0] != want {
		t.Fatalf("errors = %v, want [%q]", ve.Errors, want)
	}
}

func TestValidateStruct_BareAndDottedRule_MapDoesNotLeakUnvalidatedFields(t *testing.T) {
	p := nestedPayload{Payload: nestedModeSecret{Mode: "individual", Secret: "shouldnotleak"}}

	orderings := map[string][]Input{
		"bare_then_dotted": {
			{Name: "payload", Rules: []Rule{{Name: "required"}}},
			{Name: "payload.mode", Rules: []Rule{{Name: "required"}}},
		},
		"dotted_then_bare": {
			{Name: "payload.mode", Rules: []Rule{{Name: "required"}}},
			{Name: "payload", Rules: []Rule{{Name: "required"}}},
		},
	}

	for name, rules := range orderings {
		t.Run(name, func(t *testing.T) {
			result, err := ValidateStruct(p, rules, nil, nil)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if result.Payload.Secret != "" {
				t.Fatalf("expected Secret to be stripped, got %q", result.Payload.Secret)
			}
			if result.Payload.Mode != "individual" {
				t.Fatalf("expected Mode to survive, got %q", result.Payload.Mode)
			}
		})
	}
}

func TestValidateStruct_BareAndDottedRule_SliceDoesNotLeakUnvalidatedFields(t *testing.T) {
	p := itemsPayload{Items: []nestedModeSecret{{Mode: "individual", Secret: "shouldnotleak"}}}

	orderings := map[string][]Input{
		"bare_then_dotted": {
			{Name: "items", Rules: []Rule{{Name: "required"}}},
			{Name: "items.*.mode", Rules: []Rule{{Name: "required"}}},
		},
		"dotted_then_bare": {
			{Name: "items.*.mode", Rules: []Rule{{Name: "required"}}},
			{Name: "items", Rules: []Rule{{Name: "required"}}},
		},
	}

	for name, rules := range orderings {
		t.Run(name, func(t *testing.T) {
			result, err := ValidateStruct(p, rules, nil, nil)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if len(result.Items) != 1 {
				t.Fatalf("expected 1 item, got %d", len(result.Items))
			}
			if result.Items[0].Secret != "" {
				t.Fatalf("expected Secret to be stripped, got %q", result.Items[0].Secret)
			}
			if result.Items[0].Mode != "individual" {
				t.Fatalf("expected Mode to survive, got %q", result.Items[0].Mode)
			}
		})
	}
}

func TestValidateStruct_BareRuleWithoutDottedSibling_PassesThroughWholeObject(t *testing.T) {
	// Sin reglas hermanas con punto, una regla "bare" (ej. "payload": required)
	// debe seguir dejando pasar el objeto anidado completo, tal como lo hacía
	// antes de que se corrigiera la fuga de campos no cubiertos.
	p := nestedPayload{Payload: nestedModeSecret{Mode: "individual", Secret: "kept"}}

	rules := []Input{
		{Name: "payload", Rules: []Rule{{Name: "required"}}},
	}

	result, err := ValidateStruct(p, rules, nil, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Payload.Secret != "kept" {
		t.Fatalf("expected Secret to pass through untouched, got %q", result.Payload.Secret)
	}
}

func TestValidateStruct_DottedRuleOnly_StripsUnvalidatedFields(t *testing.T) {
	// Regresión: sin ninguna regla "bare", una regla con punto por sí sola
	// debe seguir filtrando el resto de los campos del objeto anidado.
	p := nestedPayload{Payload: nestedModeSecret{Mode: "individual", Secret: "shouldnotleak"}}

	rules := []Input{
		{Name: "payload.mode", Rules: []Rule{{Name: "required"}}},
	}

	result, err := ValidateStruct(p, rules, nil, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Payload.Secret != "" {
		t.Fatalf("expected Secret to be stripped, got %q", result.Payload.Secret)
	}
}
