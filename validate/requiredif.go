package validate

import (
	"fmt"
	"strings"
)

// getNestedValue traverses body using a dot-notation path and returns the value and whether it was found.
func getNestedValue(body map[string]any, path string) (any, bool) {
	parts := strings.Split(path, ".")
	current := any(body)

	for _, part := range parts {
		m, ok := current.(map[string]any)
		if !ok {
			return nil, false
		}
		current, ok = m[part]
		if !ok {
			return nil, false
		}
	}

	return current, true
}

// RequiredIf makes the field required when the node at opts[0] (dot-notation path) exists in body.
// If opts[1] is also provided, the field is only required when the node value equals opts[1].
func RequiredIf(inputName string, value any, body map[string]any, opts []string, sliceIndex string, allErrors map[string]any, addError func(string, string, map[string]any, string) map[string]any, customeallErrors map[string]string) map[string]any {
	if len(opts) < 1 {
		return addError(inputName, "required_if", allErrors, "La regla required_if requiere al menos 1 parámetro (ruta del nodo)")
	}

	nodePath := opts[0]
	nodeValue, exists := getNestedValue(body, nodePath)
	if !exists {
		return allErrors
	}

	if len(opts) >= 2 {
		expectedValue := opts[1]
		nodeStr := fmt.Sprintf("%v", nodeValue)
		if nodeStr != expectedValue {
			return allErrors
		}
	}

	if value == nil || value == "" {
		tmpError := fmt.Sprintf("El campo %s es requerido cuando %s está presente", inputName, nodePath)
		if len(opts) >= 2 {
			tmpError = fmt.Sprintf("El campo %s es requerido cuando %s es %s", inputName, nodePath, opts[1])
		}

		if sliceIndex != "" {
			tmpError = fmt.Sprintf("El campo en la posición %s es requerido cuando %s está presente", sliceIndex, nodePath)
		}

		tmpErrorKey := fmt.Sprintf("%s.required_if", inputName)
		if customeError, exists := customeallErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}

		return addError(inputName, "required_if", allErrors, tmpError)
	}

	return allErrors
}
