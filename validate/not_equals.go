package validate

import "fmt"

func NotEqual(input string, value any, payload map[string]any, options []string, sliceIndex string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, customeErrors map[string]string) map[string]interface{} {
	if len(options) == 0 {
		errors = addError(input, "not_equal", errors, "No se proporcionó un valor para comparar")
		return errors
	}

	if _, ok := payload[input]; !ok {
		tmpError := "El campo a comparar no existe en la carga útil"

		customeErrorKey := fmt.Sprintf("%s.not_equal", input)
		if customeError, exists := customeErrors[customeErrorKey]; exists {
			tmpError = customeError
		}

		errors = addError(input, "not_equal", errors, tmpError)
		return errors
	}

	if value == options[0] {
		tmpError := fmt.Sprintf("El valor es igual a %v", options[0])

		if sliceIndex != "" {
			tmpError = fmt.Sprintf("El valor en la posición %s es igual a %v", sliceIndex, options[0])
		}

		customeErrorKey := fmt.Sprintf("%s.not_equal", input)
		if customeError, exists := customeErrors[customeErrorKey]; exists {
			tmpError = customeError
		}

		errors = addError(input, "not_equal", errors, tmpError)
	}

	return errors
}
