package validate

import (
	"strings"
)

func RequiredWithoutAll(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}) map[string]interface{} {
	if len(options) != 1 {
		return addError(input, "required_without_all", errors, "La opción no está definida")
	}

	not_defined := true
	for _, another_input := range options {
		if _, exists_another_input := payload[another_input]; !exists_another_input {
			not_defined = false
			break
		}
	}

	// fmt.Println(input+": not_defined", not_defined)

	if _, exists_input := payload[input]; !exists_input && !not_defined {
		all_inputs := strings.Join(options, ", ")
		errors = addError(input, "required_without_all", errors, "El campo \""+input+"\" debe estar definido cuando los campos \""+all_inputs+"\" no están definidos")
		return errors
	}

	return errors
}
