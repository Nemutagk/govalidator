package validate

import "fmt"

func RequiredWithout(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, customeErrors map[string]string) map[string]interface{} {
	if len(options) != 1 {
		tmpError := "La opción no está definida"
		tmpErrorKey := fmt.Sprintf("%s.required_without", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		return addError(input, "required_without", errors, tmpError)
	}

	if _, exists := payload[options[0]]; !exists {
		if _, exists_input := payload[input]; !exists_input {
			tmpError := fmt.Sprintf("El campo \"%s\" debe estar definido cuando el campo \"%s\" no está definido", input, options[0])
			tmpErrorKey := fmt.Sprintf("%s.required_without", input)
			if customeError, exists := customeErrors[tmpErrorKey]; exists {
				tmpError = customeError
			}
			errors = addError(input, "required_without", errors, tmpError)
			return errors
		}
	}

	return errors
}
