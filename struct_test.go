package govalidator

import "testing"

type namedStatusType string

const namedStatusActive namedStatusType = "active"

type namedPriorityType int

type namedStatusPayload struct {
	Status   namedStatusType `json:"status"`
	Priority namedPriorityType `json:"priority"`
}

func TestValidateStruct_NamedStringType_PassesInRule(t *testing.T) {
	// Antes del fix, un campo de tipo con nombre (ej. "type StatusType
	// string") llegaba a las reglas envuelto en su tipo original, no como
	// string plano. La regla "in" compara con "==" contra un string literal
	// de Options, y en Go dos interfaces solo son iguales si su tipo
	// dinámico también coincide — así que "in" fallaba siempre, sin importar
	// el valor.
	p := namedStatusPayload{Status: namedStatusActive, Priority: 5}

	rules := []Input{
		{Name: "status", Rules: []Rule{{Name: "in", Options: []string{"active", "inactive"}}}},
	}

	result, err := ValidateStruct(p, rules, nil, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Status != namedStatusActive {
		t.Fatalf("got %q, want %q", result.Status, namedStatusActive)
	}
}

func TestValidateStruct_NamedIntType_PassesNumericComparisonRules(t *testing.T) {
	// Mismo problema para tipos numéricos con nombre: validate/min.go hace
	// value.(int), que fallaba para un "type Priority int" porque su tipo
	// dinámico no era exactamente "int".
	p := namedStatusPayload{Status: namedStatusActive, Priority: 5}

	rules := []Input{
		{Name: "priority", Rules: []Rule{{Name: "min", Options: []string{"1"}}, {Name: "greater_than", Options: []string{"0"}}}},
	}

	result, err := ValidateStruct(p, rules, nil, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Priority != 5 {
		t.Fatalf("got %v, want 5", result.Priority)
	}
}

type namedBoolType bool

type namedBoolPayload struct {
	Active namedBoolType `json:"active"`
}

func TestValidateStruct_NamedBoolType_PassesBooleanRule(t *testing.T) {
	// validate/boolean.go hace value.(bool), que también fallaba para un
	// tipo con nombre cuyo subyacente es bool.
	p := namedBoolPayload{Active: true}

	rules := []Input{
		{Name: "active", Rules: []Rule{{Name: "boolean"}}},
	}

	if _, err := ValidateStruct(p, rules, nil, nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateStruct_PlainPrimitiveTypes_StillWork(t *testing.T) {
	// Regresión: campos con tipos primitivos sin nombre (el caso común) no
	// deben cambiar de comportamiento con el fix.
	p := statusPayload{Status: "individual"}

	rules := []Input{
		{Name: "status", Rules: []Rule{{Name: "in", Options: []string{"individual", "batch"}}}},
	}

	result, err := ValidateStruct(p, rules, nil, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Status != "individual" {
		t.Fatalf("got %q, want %q", result.Status, "individual")
	}
}
