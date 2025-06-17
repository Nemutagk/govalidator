package validate

import "strings"

func RequiredWithAll(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}) map[string]interface{} {
	if len(options) != 1 {
		return addError(input, "required_with_all", errors, "La opción no está definida")
	}

	all_defined := true
	for _, another_input := range options {
		if _, exists_another_input := payload[another_input]; !exists_another_input {
			all_defined = false
			break
		}
	}

	if _, exists_input := payload[input]; !exists_input && all_defined {
		all_inputs := strings.Join(options, ", ")
		errors = addError(input, "required_with_all", errors, "El campo \""+input+"\" debe estar definido cuando los campos \""+all_inputs+"\" están definidos")
		return errors
	}

	return errors
}
