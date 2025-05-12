package validate

import "strconv"

func Max(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}) map[string]interface{} {
	value := payload[input]

	max, err := strconv.ParseInt(options[0], 10, 64)
	if err != nil {
		errors = addError(input, "min", errors, "The field "+input+" must be a number")
		return errors
	}

	if _, ok := value.(string); ok {
		strlen := len(value.(string))

		if strlen > int(max) {
			errors = addError(input, "min", errors, "The field "+input+" must be at most "+options[0]+" characters long")
		}
	}

	if _, ok := value.(int); ok {
		intValue := value.(int)
		if intValue > int(max) {
			errors = addError(input, "min", errors, "The field "+input+" must be at most "+options[0]+"")
		}
	}

	if _, ok := value.(float64); ok {
		floatValue := value.(float64)
		if floatValue > float64(max) {
			errors = addError(input, "max", errors, "The field "+input+" must be at most "+options[0]+"")
		}
	}

	return errors
}
