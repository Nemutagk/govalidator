package validate

import (
	"fmt"
	"strings"
)

// RequiredIfAll makes the field required only when ALL key/value pairs in opts hold true (AND).
// opts must be a flat, even-length slice of pairs: [path1, value1, path2, value2, ...].
func RequiredIfAll(inputName string, value any, body map[string]any, opts []string, sliceIndex string, allErrors map[string]any, addError func(string, string, map[string]any, string) map[string]any, customeallErrors map[string]string) (map[string]any, bool) {
	total := len(opts)
	if total < 2 || total%2 != 0 {
		return addError(inputName, "required_if_all", allErrors, "La regla required_if_all requiere un número par de parámetros (pares de ruta del nodo y valor esperado)"), true
	}

	pairs := total / 2
	conditions := make([]string, 0, pairs)
	for i := 0; i < pairs; i++ {
		nodePath, expected := opts[i*2], opts[i*2+1]
		nodeValue, exists := getNestedValue(body, nodePath)
		if !exists || fmt.Sprintf("%v", nodeValue) != expected {
			return allErrors, false
		}
		conditions = append(conditions, fmt.Sprintf("%s es %s", nodePath, expected))
	}

	if value == nil || value == "" {
		conditionsStr := strings.Join(conditions, " y ")
		tmpError := fmt.Sprintf("El campo %s es requerido cuando %s", inputName, conditionsStr)
		if sliceIndex != "" {
			tmpError = fmt.Sprintf("El campo en la posición %s es requerido cuando %s", sliceIndex, conditionsStr)
		}

		tmpErrorKey := fmt.Sprintf("%s.required_if_all", inputName)
		if customeError, exists := customeallErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		return addError(inputName, "required_if_all", allErrors, tmpError), true
	}

	return allErrors, false
}
