package validate

import (
	"fmt"
	"time"
)

func DateFormat(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, customeErrors map[string]string) map[string]interface{} {
	if len(options) == 0 {
		errors = addError(input, "date_format", errors, "El formato de fecha no está definido")
		return errors
	}

	value := payload[input]

	if value == nil || value == "" {
		return errors
	}

	formato := options[0]

	_, err_date := time.Parse(formato, value.(string))

	if err_date != nil {
		tmpError := "El formato de la fecha es inválido"

		customeErrorKey := fmt.Sprintf("%s.date_format", input)
		if customeError, exists := customeErrors[customeErrorKey]; exists {
			tmpError = customeError
		}

		errors = addError(input, "date_format", errors, tmpError)
	}

	return errors
}
