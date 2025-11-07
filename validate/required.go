package validate

import "fmt"

func Required(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, customeErrors map[string]string) map[string]interface{} {
	if _, exists_input := payload[input]; !exists_input {
		tmpError := fmt.Sprintf("El campo %s no está definido", input)
		tmpErrorKey := fmt.Sprintf("%s.required", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		errors = addError(input, "required", errors, tmpError)
		return errors
	}

	if payload[input] == "" || payload[input] == nil {
		tmpError := fmt.Sprintf("El campo %s está vacío", input)
		tmpErrorKey := fmt.Sprintf("%s.required", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		errors = addError(input, "required", errors, tmpError)
		return errors
	}

	return errors
}
