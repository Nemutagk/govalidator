package validate

import (
	"fmt"
	"regexp"
)

func Email(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, customeErrors map[string]string) map[string]interface{} {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	value, ok := payload[input].(string)
	if ok && !emailRegex.MatchString(value) {
		tmpError := "El campo no es un correo electrónico válido"

		customeErrorKey := fmt.Sprintf("%s.email", input)
		if customeError, exists := customeErrors[customeErrorKey]; exists {
			tmpError = customeError
		}

		errors = addError(input, "email", errors, tmpError)
	}

	return errors
}
