package validate

import "fmt"

// InIf valida que el valor esté en la lista de valores permitidos (opts[2:])
// solo cuando el nodo en opts[0] (ruta con notación de punto sobre el body
// raíz) es igual a opts[1]. Si la condición no se cumple (el nodo no existe o
// tiene otro valor), la regla no hace nada — igual que required_if.
func InIf(input string, value any, body map[string]any, opts []string, sliceIndex string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, customeErrors map[string]string) map[string]interface{} {
	if len(opts) < 3 {
		return addError(input, "in_if", errors, "La regla in_if requiere al menos 3 parámetros (ruta del nodo, valor esperado y al menos un valor permitido)")
	}

	nodePath, expected := opts[0], opts[1]
	allowed := opts[2:]

	nodeValue, exists := getNestedValue(body, nodePath)
	if !exists {
		return errors
	}

	if fmt.Sprintf("%v", nodeValue) != expected {
		return errors
	}

	if value == nil || value == "" {
		return errors
	}

	for _, option := range allowed {
		if option == value {
			return errors
		}
	}

	tmpError := "No se encontró el valor en las opciones permitidas"

	if sliceIndex != "" {
		tmpError = fmt.Sprintf("El valor en la posición %s no se encontró en las opciones permitidas", sliceIndex)
	}

	customeErrorKey := fmt.Sprintf("%s.in_if", input)
	if customeError, exists := customeErrors[customeErrorKey]; exists {
		tmpError = customeError
	}

	return addError(input, "in_if", errors, tmpError)
}
