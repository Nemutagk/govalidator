package validate

import (
	"fmt"
	"regexp"
)

func Email(input string, value any, payload map[string]any, options []string, sliceIndex string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, customeErrors map[string]string) map[string]interface{} {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	value, ok := payload[input].(string)
	if ok && !emailRegex.MatchString(value.(string)) {
		tmpError := "El campo no es un correo electrónico válido"

		if sliceIndex != "" {
			tmpError = fmt.Sprintf("El campo en la posición %s no es un correo electrónico válido", sliceIndex)
		}

		customeErrorKey := fmt.Sprintf("%s.email", input)
		if customeError, exists := customeErrors[customeErrorKey]; exists {
			tmpError = customeError
		}

		errors = addError(input, "email", errors, tmpError)
	}

	return errors
}
