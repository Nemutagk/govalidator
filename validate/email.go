package validate

import "regexp"

func Email(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}) map[string]interface{} {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	value, ok := payload[input].(string)
	if ok && !emailRegex.MatchString(value) {
		errors = addError(input, "email", errors, "El valor no es un correo electrónico válido")
	}

	return errors
}
