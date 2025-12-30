package validate

import (
	"fmt"
)

func RequiredWith(input string, value any, payload map[string]any, options []string, sliceIndex string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, customeErrors map[string]string) map[string]interface{} {
	if len(options) < 1 && len(options) > 2 {
		tmpError := "La opción no esta definida"

		if sliceIndex != "" {
			tmpError = fmt.Sprintf("La opción en la posición %s no esta definida", sliceIndex)
		}

		tmpErrorKey := fmt.Sprintf("%s.required_with", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		return addError(input, "required_with", errors, tmpError)
	}

	existsWithValue, exists_input := payload[options[0]]
	if !exists_input && len(options) == 1 {
		tmpError := fmt.Sprintf("El campo '%s' debe estar definido cuando el campo '%s' está definido", input, options[0])

		if sliceIndex != "" {
			tmpError = fmt.Sprintf("El campo '%s' en la posición %s debe estar definido cuando el campo '%s' está definido", input, sliceIndex, options[0])
		}

		tmpErrorKey := fmt.Sprintf("%s.required_with", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		errors = addError(input, "required_with", errors, tmpError)
		return errors
	}

	if !exists_input && len(options) == 2 && existsWithValue != options[1] {
		tmpError := fmt.Sprintf("El campo '%s' debe estar definido cuando el campo '%s' está definido y el valor es '%s'", input, options[0], options[1])

		if sliceIndex != "" {
			tmpError = fmt.Sprintf("El campo '%s' en la posición %s debe estar definido cuando el campo '%s' está definido y el valor es '%s'", input, sliceIndex, options[0], options[1])
		}

		tmpErrorKey := fmt.Sprintf("%s.required_with", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		errors = addError(input, "required_with", errors, tmpError)
		return errors
	}

	return errors
}
