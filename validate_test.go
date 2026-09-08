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

type statusPayload struct {
	Status string `json:"status"`
}

type confirmablePayload struct {
	Email             string `json:"email"`
	EmailConfirmation string `json:"email_confirmation"`
}

func TestValidateStruct_NormalizerLowerMakesInRulePass(t *testing.T) {
	p := statusPayload{Status: "INDIVIDUAL"}

	rules := []Input{
		{
			Name:        "status",
			Normalizers: []Normalizer{{Name: "lower"}},
			Rules:       []Rule{{Name: "in", Options: []string{"individual", "batch"}}},
		},
	}

	result, err := ValidateStruct(p, rules, nil, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Status != "individual" {
		t.Fatalf("got %q, want %q", result.Status, "individual")
	}
}

func TestValidateRequest_NormalizerMutatesCallerMapInPlace(t *testing.T) {
	// Comportamiento acordado: los normalizers mutan en sitio el body que se
	// les pasa, para que otras reglas que lean ese mismo campo vean el valor
	// ya normalizado. Si se llama ValidateRequest directo con un mapa, ese
	// mapa del caller queda modificado.
	body := map[string]any{"status": "  ACTIVE  "}

	rules := []Input{
		{
			Name:        "status",
			Normalizers: []Normalizer{{Name: "trim"}, {Name: "lower"}},
			Rules:       []Rule{{Name: "in", Options: []string{"active", "inactive"}}},
		},
	}

	if _, err := ValidateRequest(body, rules, nil, nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if body["status"] != "active" {
		t.Fatalf("expected caller's map to be mutated in place, got %q", body["status"])
	}
}

func TestValidateStruct_NormalizerOnArrayDottedElement(t *testing.T) {
	// Blinda el fix en el reconstructor de elemInput (validate.go): sin él,
	// un Normalizer declarado sobre "items.*.mode" se perdería en cada
	// elemento del arreglo.
	p := itemsPayload{Items: []nestedModeSecret{{Mode: "individual"}}}

	rules := []Input{
		{
			Name:        "items.*.mode",
			Normalizers: []Normalizer{{Name: "upper"}},
			Rules:       []Rule{{Name: "in", Options: []string{"INDIVIDUAL", "BATCH"}}},
		},
	}

	result, err := ValidateStruct(p, rules, nil, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Items) != 1 || result.Items[0].Mode != "INDIVIDUAL" {
		t.Fatalf("unexpected result: %+v", result.Items)
	}
}

func TestValidateStruct_NormalizerVisibleToConfirmationRule(t *testing.T) {
	// Prueba real de por qué el normalizer tiene que escribir de vuelta en
	// body y no solo en la variable local: "confirmation" lee el campo
	// directo del body, no el value que le llega a "email".
	p := confirmablePayload{Email: "  a@b.com  ", EmailConfirmation: "a@b.com"}

	rules := []Input{
		{
			Name:        "email",
			Normalizers: []Normalizer{{Name: "trim"}},
			Rules:       []Rule{{Name: "confirmation"}},
		},
	}

	if _, err := ValidateStruct(p, rules, nil, nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateStruct_UnknownNormalizerName_ProducesError(t *testing.T) {
	p := statusPayload{Status: "individual"}

	rules := []Input{
		{Name: "status", Normalizers: []Normalizer{{Name: "does_not_exist"}}},
	}

	_, err := ValidateStruct(p, rules, nil, nil)
	ve, ok := AsValidationError(err)
	if !ok {
		t.Fatalf("expected a ValidationError, got %v", err)
	}
	want := "status: El normalizador does_not_exist no es válido"
	if len(ve.Errors) != 1 || ve.Errors[0] != want {
		t.Fatalf("errors = %v, want [%q]", ve.Errors, want)
	}
}
