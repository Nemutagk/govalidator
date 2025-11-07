package validate

import "fmt"

func RequiredWith(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, customeErrors map[string]string) map[string]interface{} {
	if len(options) != 1 {
		tmpError := "La opción no esta definida"
		tmpErrorKey := fmt.Sprintf("%s.required_with", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		return addError(input, "required_with", errors, tmpError)
	}

	if _, exists_input := payload[input]; !exists_input {
		tmpError := fmt.Sprintf("El campo \"%s\" debe estar definido cuando el campo \"%s\" está definido", input, options[0])
		tmpErrorKey := fmt.Sprintf("%s.required_with", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		errors = addError(input, "required_with", errors, tmpError)
		return errors
	}

	return errors
}
