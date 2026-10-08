package validate

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

var uuidFormats = map[string]bool{"canonical": true, "simple": true, "braces": true, "urn": true}

var uuidVersions = map[string]uuid.Version{"4": 4, "5": 5, "6": 6, "7": 7, "8": 8}

func UUID(input string, value any, payload map[string]any, options []string, sliceIndex string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, customeErrors map[string]string) map[string]interface{} {
	fail := func(tmpError string, tmpErrorWithIndex string) map[string]interface{} {
		if sliceIndex != "" {
			tmpError = tmpErrorWithIndex
		}

		tmpErrorKey := fmt.Sprintf("%s.uuid", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		return addError(input, "uuid", errors, tmpError)
	}

	formats := []string{}
	versions := map[uuid.Version]bool{}
	invalidOption := false
	for _, opt := range options {
		if uuidFormats[opt] {
			formats = append(formats, opt)
		} else if v, ok := uuidVersions[opt]; ok {
			versions[v] = true
		} else {
			invalidOption = true
			errors = fail(
				fmt.Sprintf("La opción %s de la regla uuid del campo %s no es válida", opt, input),
				fmt.Sprintf("La opción %s de la regla uuid del campo %s en la posición %s no es válida", opt, input, sliceIndex),
			)
		}
	}
	if invalidOption {
		return errors
	}

	str, ok := value.(string)
	if !ok || str == "" {
		return errors
	}

	if len(formats) == 0 {
		formats = []string{"canonical"}
	}

	for _, format := range formats {
		canonical, ok := uuidToCanonical(format, str)
		if !ok {
			continue
		}

		parsed, err := uuid.Parse(canonical)
		if err != nil {
			continue
		}

		if parsed.Variant() != uuid.RFC4122 {
			continue
		}

		version := parsed.Version()
		if version < 4 || version > 8 {
			continue
		}
		if len(versions) > 0 && !versions[version] {
			continue
		}

		if uuidIsRepetitive(canonical) {
			continue
		}

		return errors
	}

	return fail(
		fmt.Sprintf("El campo %s no es un UUID válido", input),
		fmt.Sprintf("El campo %s en la posición %s no es un UUID válido", input, sliceIndex),
	)
}

func uuidToCanonical(format string, s string) (string, bool) {
	switch format {
	case "canonical":
		return s, len(s) == 36
	case "simple":
		if len(s) != 32 {
			return "", false
		}
		return s[0:8] + "-" + s[8:12] + "-" + s[12:16] + "-" + s[16:20] + "-" + s[20:], true
	case "braces":
		if len(s) != 38 || s[0] != '{' || s[37] != '}' {
			return "", false
		}
		return s[1:37], true
	case "urn":
		const prefix = "urn:uuid:"
		if len(s) != len(prefix)+36 || !strings.EqualFold(s[:len(prefix)], prefix) {
			return "", false
		}
		return s[len(prefix):], true
	}

	return "", false
}

// Un UUID cuyos 30 dígitos hex (sin versión ni variante) son todos iguales es inventado.
func uuidIsRepetitive(canonical string) bool {
	digits := strings.ReplaceAll(strings.ToLower(canonical), "-", "")
	digits = digits[:12] + digits[13:16] + digits[17:]

	return strings.Count(digits, digits[:1]) == len(digits)
}
