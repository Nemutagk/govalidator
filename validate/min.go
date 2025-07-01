package validate

import "strconv"

func Min(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}) map[string]interface{} {
	value := payload[input]

	min, err := strconv.ParseInt(options[0], 10, 64)
	if err != nil {
		errors = addError(input, "min", errors, "El campo "+input+" debe ser un número")
		return errors
	}

	if _, ok := value.(string); ok {
		strlen := len(value.(string))

		if strlen < int(min) {
			errors = addError(input, "min", errors, "El campo "+input+" debe tener al menos "+options[0]+" caracteres")
		}
	}

	if _, ok := value.(int); ok {
		intValue := value.(int)
		if intValue < int(min) {
			errors = addError(input, "min", errors, "El campo "+input+" debe ser al menos "+options[0]+"")
		}
	}

	if _, ok := value.(float64); ok {
		floatValue := value.(float64)
		if floatValue < float64(min) {
			errors = addError(input, "min", errors, "El campo "+input+" debe ser al menos "+options[0]+"")
		}
	}

	return errors
}
