package validate

import "regexp"

func Password(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}) map[string]interface{} {
	if _, exists_input := payload[input]; !exists_input {
		return errors
	}

	value, ok := payload[input].(string)
	if !ok {
		errors = addError("password", "type", errors, "Password must be a string")
		return errors
	}

	if len(value) < 6 {
		errors = addError("password", "min:6", errors, "Password must be at least 6 characters long")
	}

	if match, _ := regexp.MatchString("[0-9]", value); !match {
		errors = addError("password", "regex", errors, "Password must contain at least one number")
	}

	if match, _ := regexp.MatchString("[a-z]", value); !match {
		errors = addError("password", "regex", errors, "Password must contain at least one lowercase letter")
	}

	if match, _ := regexp.MatchString("[A-Z]", value); !match {
		errors = addError("password", "regex", errors, "Password must contain at least one uppercase letter")
	}

	if match, _ := regexp.MatchString(`[#$%&/()!_-]+`, value); !match {
		errors = addError("password", "regex", errors, "Password must contain at least one special character ($#%&/()!_-)")
	}

	return errors
}
