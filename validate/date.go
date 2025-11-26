package validate

import (
	"fmt"
	"time"
)

func Date(input string, value any, payload map[string]any, options []string, sliceIndex string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, customeErrors map[string]string) map[string]interface{} {
	if value == nil || value == "" {
		return errors
	}

	layout := "2006-01-02T15:04:05"
	if len(options) > 1 && options[0] != "" {
		layout = options[0]
	}

	_, exists := payload[input].(string)
	if !exists {
		tmpError := fmt.Sprintf("El campo \"%s\" debe ser una fecha válida con el formato \"%s\"", input, layout)

		if sliceIndex != "" {
			tmpError = fmt.Sprintf("El campo \"%s\" en la posición %s debe ser una fecha válida con el formato \"%s\"", input, sliceIndex, layout)
		}

		tmpErrorKey := fmt.Sprintf("%s.date", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		errors = addError(input, "date", errors, tmpError)
		return errors
	}

	if _, err := time.Parse(layout, value.(string)); err != nil {
		tmpError := fmt.Sprintf("El campo \"%s\" debe ser una fecha válida con el formato \"%s\"", input, layout)

		if sliceIndex != "" {
			tmpError = fmt.Sprintf("El campo \"%s\" en la posición %s debe ser una fecha válida con el formato \"%s\"", input, sliceIndex, layout)
		}

		tmpErrorKey := fmt.Sprintf("%s.date", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		errors = addError(input, "date", errors, tmpError)
	}

	return errors
}
