package validate

import "github.com/Nemutagk/govalidator/helper"

func In(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}) map[string]interface{} {
	value := payload[input]

	if value == nil || value == "" {
		return errors
	}

	stringValue, ok := value.(string)

	if !ok {
		errors = addError(input, "in", errors, "The value is not a string")
		return errors
	}

	if !helper.SliceContains(options, stringValue) {
		errors = addError(input, "in", errors, "The value is invalid")
	}

	return errors
}
