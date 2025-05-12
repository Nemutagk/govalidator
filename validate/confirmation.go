package validate

import "fmt"

func Confirmation(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}) map[string]interface{} {
	if _, exists_input := payload[input]; !exists_input {
		fmt.Println("validate confirmation:input password does not exists")
		return errors
	}

	if payload[input] == "" {
		errors = addError(input, "confirmation", errors, "The input "+input+" is empty")
		return errors
	}

	if val, exists_confirmation := payload["password_confirmation"]; !exists_confirmation && (val == nil || val == "" || val == false) {
		errors = addError("password_confirmation", "confirmation", errors, "The input password_confirmation is not defined")
	}

	if payload[input] != payload["password_confirmation"] {
		errors = addError(input, "confirmation", errors, "The password is not same for confirmation")
	}

	return errors
}
