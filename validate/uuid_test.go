package validate

import (
	"testing"

	"github.com/google/uuid"
)

const (
	uuidV1  = "c232ab00-9414-11ec-b3c8-9f6bdeced846"
	uuidV3  = "6fa459ea-ee8a-3ca4-894e-db77e160355e"
	uuidV4  = "550e8400-e29b-41d4-a716-446655440000"
	uuidV5  = "886313e1-3b8a-5372-9b90-0c9aee199e5d"
	uuidV6  = "1ec9414c-232a-6b00-b3c8-9f6bdeced846"
	uuidV7  = "018f3c5e-7b1a-7c3d-9e2f-4a5b6c7d8e9f"
	uuidV8  = "2489e9ad-2ee2-8e00-8ec9-32d5f69181c0"
	uuidNil = "00000000-0000-0000-0000-000000000000"
	uuidMax = "ffffffff-ffff-ffff-ffff-ffffffffffff"
)

func runUUID(value any, sliceIndex string, customErrors map[string]string, options ...string) map[string]interface{} {
	errors := make(map[string]interface{})
	return UUID("id", value, map[string]any{}, options, sliceIndex, errors, testAddError, customErrors)
}

func assertUUIDValid(t *testing.T, value any, options ...string) {
	t.Helper()
	if errors := runUUID(value, "", map[string]string{}, options...); len(errors) != 0 {
		t.Fatalf("value %v options %v: expected no errors, got %v", value, options, errors)
	}
}

func assertUUIDInvalid(t *testing.T, value any, options ...string) {
	t.Helper()
	errors := runUUID(value, "", map[string]string{}, options...)
	msgs, found := getErrorMsgs(errors, "id", "uuid")
	want := "El campo id no es un UUID válido"
	if !found || msgs[0] != want {
		t.Fatalf("value %v options %v: msgs = %v, want %q", value, options, msgs, want)
	}
}

func TestUUID_DefaultAcceptsV4ToV8(t *testing.T) {
	for _, v := range []string{uuidV4, uuidV5, uuidV6, uuidV7, uuidV8} {
		assertUUIDValid(t, v)
	}
}

func TestUUID_DefaultAcceptsGeneratedV4AndV7(t *testing.T) {
	v4, err := uuid.NewRandom()
	if err != nil {
		t.Fatal(err)
	}
	v7, err := uuid.NewV7()
	if err != nil {
		t.Fatal(err)
	}

	assertUUIDValid(t, v4.String())
	assertUUIDValid(t, v7.String())
}

func TestUUID_RejectsOldVersions(t *testing.T) {
	assertUUIDInvalid(t, uuidV1)
	assertUUIDInvalid(t, uuidV3)
}

func TestUUID_RejectsNilAndMax(t *testing.T) {
	assertUUIDInvalid(t, uuidNil)
	assertUUIDInvalid(t, uuidMax)
}

func TestUUID_RejectsInvalidVariant(t *testing.T) {
	assertUUIDInvalid(t, "550e8400-e29b-41d4-c716-446655440000")
}

func TestUUID_RejectsRepetitiveInvented(t *testing.T) {
	assertUUIDInvalid(t, "11111111-1111-4111-8111-111111111111")
	assertUUIDInvalid(t, "00000000-0000-4000-8000-000000000000")
	assertUUIDInvalid(t, "aaaaaaaa-aaaa-7aaa-9aaa-aaaaaaaaaaaa")
}

func TestUUID_RejectsBadFormat(t *testing.T) {
	assertUUIDInvalid(t, "no-es-un-uuid")
	assertUUIDInvalid(t, "550e8400-e29b-41d4-a716-44665544000")
	assertUUIDInvalid(t, "550e8400-e29b-41d4-a716-44665544000g")
	assertUUIDInvalid(t, "550e8400e29b41d4a716446655440000")
	assertUUIDInvalid(t, "{"+uuidV4+"}")
	assertUUIDInvalid(t, "urn:uuid:"+uuidV4)
}

func TestUUID_AcceptsUppercaseHex(t *testing.T) {
	assertUUIDValid(t, "550E8400-E29B-41D4-A716-446655440000")
}

func TestUUID_IgnoresEmptyNilAndNonString(t *testing.T) {
	for _, v := range []any{"", nil, 5, 3.14, true, []string{"x"}} {
		assertUUIDValid(t, v)
	}
}

func TestUUID_FormatSimple(t *testing.T) {
	assertUUIDValid(t, "550e8400e29b41d4a716446655440000", "simple")
	assertUUIDInvalid(t, uuidV4, "simple")
	assertUUIDInvalid(t, "550e8400e29b41d4a71644665544000g", "simple")
}

func TestUUID_FormatBraces(t *testing.T) {
	assertUUIDValid(t, "{"+uuidV4+"}", "braces")
	assertUUIDInvalid(t, uuidV4, "braces")
	assertUUIDInvalid(t, "{"+uuidV4+")", "braces")
}

func TestUUID_FormatUrn(t *testing.T) {
	assertUUIDValid(t, "urn:uuid:"+uuidV4, "urn")
	assertUUIDValid(t, "URN:UUID:"+uuidV4, "urn")
	assertUUIDInvalid(t, uuidV4, "urn")
	assertUUIDInvalid(t, "urn:guid:"+uuidV4, "urn")
}

func TestUUID_FormatCanonicalExplicit(t *testing.T) {
	assertUUIDValid(t, uuidV4, "canonical")
	assertUUIDInvalid(t, "{"+uuidV4+"}", "canonical")
}

func TestUUID_MultipleFormatsAreOR(t *testing.T) {
	assertUUIDValid(t, uuidV4, "canonical", "urn")
	assertUUIDValid(t, "urn:uuid:"+uuidV4, "canonical", "urn")
	assertUUIDInvalid(t, "{"+uuidV4+"}", "canonical", "urn")
}

func TestUUID_SpecificVersion(t *testing.T) {
	assertUUIDValid(t, uuidV7, "7")
	assertUUIDInvalid(t, uuidV4, "7")
}

func TestUUID_MultipleVersionsAreOR(t *testing.T) {
	assertUUIDValid(t, uuidV4, "7", "4")
	assertUUIDValid(t, uuidV7, "4", "7")
	assertUUIDInvalid(t, uuidV5, "4", "7")
}

func TestUUID_FormatAndVersionAreAND(t *testing.T) {
	assertUUIDValid(t, "urn:uuid:"+uuidV7, "urn", "7")
	assertUUIDValid(t, "urn:uuid:"+uuidV7, "7", "urn")
	assertUUIDInvalid(t, uuidV7, "urn", "7")
	assertUUIDInvalid(t, "urn:uuid:"+uuidV4, "urn", "7")
}

func TestUUID_OptionOrderDoesNotMatter(t *testing.T) {
	assertUUIDValid(t, "{"+uuidV7+"}", "7", "4", "braces", "urn")
	assertUUIDValid(t, "{"+uuidV7+"}", "braces", "urn", "4", "7")
}

func TestUUID_InvalidOption(t *testing.T) {
	for _, opt := range []string{"xyz", "1", "2", "3", "9", ""} {
		errors := runUUID(uuidV4, "", map[string]string{}, opt)

		msgs, found := getErrorMsgs(errors, "id", "uuid")
		want := "La opción " + opt + " de la regla uuid del campo id no es válida"
		if !found || msgs[0] != want {
			t.Fatalf("option %q: msgs = %v, want %q", opt, msgs, want)
		}
	}
}

func TestUUID_InvalidOption_WithSliceIndex(t *testing.T) {
	errors := runUUID(uuidV4, "1", map[string]string{}, "xyz")

	msgs, found := getErrorMsgs(errors, "id", "uuid")
	want := "La opción xyz de la regla uuid del campo id en la posición 1 no es válida"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestUUID_InvalidOption_CustomError(t *testing.T) {
	errors := runUUID(uuidV4, "", map[string]string{"id.uuid": "mensaje personalizado"}, "xyz")

	msgs, found := getErrorMsgs(errors, "id", "uuid")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}

func TestUUID_Invalid_WithSliceIndex(t *testing.T) {
	errors := runUUID("nope", "2", map[string]string{})

	msgs, found := getErrorMsgs(errors, "id", "uuid")
	want := "El campo id en la posición 2 no es un UUID válido"
	if !found || msgs[0] != want {
		t.Fatalf("msgs = %v, want %q", msgs, want)
	}
}

func TestUUID_Invalid_CustomError(t *testing.T) {
	errors := runUUID("nope", "", map[string]string{"id.uuid": "mensaje personalizado"})

	msgs, found := getErrorMsgs(errors, "id", "uuid")
	if !found || msgs[0] != "mensaje personalizado" {
		t.Fatalf("expected custom error message, got %v", errors)
	}
}
