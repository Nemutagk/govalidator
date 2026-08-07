package validate

import "testing"

func TestIp_ValidIPv4Passes(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Ip("ip", "192.168.1.1", map[string]any{}, []string{}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestIp_ValidIPv6Passes(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Ip("ip", "2001:db8::1", map[string]any{}, []string{}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestIp_NilValueNoError(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Ip("ip", nil, map[string]any{}, []string{}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors for nil value, got %v", errors)
	}
}

func TestIp_EmptyStringNoError(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Ip("ip", "", map[string]any{}, []string{}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors for empty string, got %v", errors)
	}
}

func TestIp_TrimsWhitespaceSingle(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Ip("ip", " 127.0.0.1 ", map[string]any{}, []string{}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors for whitespace-padded valid ip, got %v", errors)
	}
}

func TestIp_InvalidSingleFails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Ip("ip", "badip", map[string]any{}, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "ip", "ip")
	want := "La dirección IP badip no es una dirección IP válida"
	if !found || len(msgs) != 1 || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestIp_InvalidSingleWithSliceIndex(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Ip("ip", "badip", map[string]any{}, []string{}, "3", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "ip", "ip")
	want := "La dirección IP badip en la posición 3 no es válida"
	if !found || len(msgs) != 1 || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestIp_InvalidSingleCustomError(t *testing.T) {
	errors := make(map[string]interface{})
	customErrors := map[string]string{"ip.ip": "mensaje personalizado"}
	errors = Ip("ip", "badip", map[string]any{}, []string{}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "ip", "ip")
	if !found || len(msgs) != 1 || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestIp_ListAllValidPasses(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Ip("ip", "192.168.1.1,10.0.0.1", map[string]any{}, []string{}, "", errors, testAddError, map[string]string{})

	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %v", errors)
	}
}

func TestIp_ListOneInvalidFails(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Ip("ip", "192.168.1.1,badip", map[string]any{}, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "ip", "ip")
	want := "La dirección IP 1:badip no es válida"
	if !found || len(msgs) != 1 || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestIp_ListMultipleInvalidAccumulatesErrors(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Ip("ip", "badip,alsobad", map[string]any{}, []string{}, "", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "ip", "ip")
	if !found || len(msgs) != 2 {
		t.Fatalf("expected 2 error messages, got %v", errors)
	}
	if msgs[0] != "La dirección IP 0:badip no es válida" {
		t.Fatalf("msgs[0] = %q, unexpected", msgs[0])
	}
	if msgs[1] != "La dirección IP 1:alsobad no es válida" {
		t.Fatalf("msgs[1] = %q, unexpected", msgs[1])
	}
}

func TestIp_ListWithSliceIndexUsesPositionFormat(t *testing.T) {
	errors := make(map[string]interface{})
	errors = Ip("ip", "192.168.1.1,badip", map[string]any{}, []string{}, "2", errors, testAddError, map[string]string{})

	msgs, found := getErrorMsgs(errors, "ip", "ip")
	want := "La dirección IP badip en la posición 2 no es válida"
	if !found || len(msgs) != 1 || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestIp_ListCustomError(t *testing.T) {
	errors := make(map[string]interface{})
	customErrors := map[string]string{"ip.ip": "mensaje personalizado"}
	errors = Ip("ip", "192.168.1.1,badip", map[string]any{}, []string{}, "", errors, testAddError, customErrors)

	msgs, found := getErrorMsgs(errors, "ip", "ip")
	if !found || len(msgs) != 1 || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}
