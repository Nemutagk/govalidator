package validate

import "fmt"

func In(input string, value any, payload map[string]any, options []string, sliceIndex string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, customeErrors map[string]string) map[string]interface{} {
	if value == nil || value == "" {
		return errors
	}

	encontrado := false
	for _, option := range options {
		if option == value {
			encontrado = true
			break
		}
	}

	if !encontrado {
		tmpError := "No se encontró el valor en las opciones permitidas"

		if sliceIndex != "" {
			tmpError = fmt.Sprintf("El valor en la posición %s no se encontró en las opciones permitidas", sliceIndex)
		}

		customeErrorKey := fmt.Sprintf("%s.in", input)
		if customeError, exists := customeErrors[customeErrorKey]; exists {
			tmpError = customeError
		}

		errors = addError(input, "in", errors, tmpError)
	}

	return errors
}
