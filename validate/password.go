package validate

import "regexp"

func Password(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}) map[string]interface{} {
	if _, exists_input := payload[input]; !exists_input {
		return errors
	}

	value, ok := payload[input].(string)
	if !ok {
		errors = addError("password", "type", errors, "La contraseña debe ser una cadena de texto")
		return errors
	}

	if len(value) < 6 {
		errors = addError("password", "min:6", errors, "La contraseña debe tener al menos 6 caracteres")
	}

	if match, _ := regexp.MatchString("[0-9]", value); !match {
		errors = addError("password", "regex", errors, "La contraseña debe contener al menos un número")
	}

	if match, _ := regexp.MatchString("[a-z]", value); !match {
		errors = addError("password", "regex", errors, "La contraseña debe contener al menos una letra minúscula")
	}

	if match, _ := regexp.MatchString("[A-Z]", value); !match {
		errors = addError("password", "regex", errors, "La contraseña debe contener al menos una letra mayúscula")
	}

	if match, _ := regexp.MatchString(`[#$%&/()!_-]+`, value); !match {
		errors = addError("password", "regex", errors, "La contraseña debe contener al menos un carácter especial ($#%&/()!_-)")
	}

	return errors
}
