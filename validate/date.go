package validate

import "time"

func Date(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}) map[string]interface{} {
	value := payload[input]

	if value == nil || value == "" {
		return errors
	}

	layout := "2006-01-02T15:04:05"
	if len(options) > 1 && options[0] != "" {
		layout = options[0]
	}

	_, exists := payload[input].(string)
	if !exists {
		errors = addError(input, "date", errors, "El campo "+input+" debe ser una fecha válida con el formato "+layout)
		return errors
	}

	if _, err := time.Parse(layout, value.(string)); err != nil {
		errors = addError(input, "date", errors, "El campo "+input+" debe ser una fecha válida con el formato "+layout)
	}

	return errors
}
