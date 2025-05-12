package validate

import "strings"

func RequiredWithout(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}) map[string]interface{} {
	if len(options) != 1 {
		return addError(input, "required_without", errors, "The options is not defined")
	}

	not_defined := true
	for _, another_input := range options {
		if _, exists_another_input := payload[another_input]; !exists_another_input {
			not_defined = false
			break
		}
	}

	if _, exists_input := payload[input]; !exists_input && !not_defined {
		all_inputs := strings.Join(options, ", ")
		errors = addError(input, "required_without", errors, "The field \""+input+"\" must be defined when any of the fields \""+all_inputs+"\" is not defined")
		return errors
	}

	return errors
}
