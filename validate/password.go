package validate

import (
	"fmt"
	"regexp"
)

func Password(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, customeErrors map[string]string) map[string]interface{} {
	if _, exists_input := payload[input]; !exists_input {
		return errors
	}

	value, ok := payload[input].(string)
	if !ok {
		tmpError := "La contraseña debe ser una cadena de texto"
		tmpErrorKey := fmt.Sprintf("%v.password", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		errors = addError("password", "type", errors, tmpError)
		return errors
	}

	if len(value) < 6 {
		tmpError := "La contraseña debe tener al menos 6 caracteres"
		tmpErrorKey := fmt.Sprintf("%s.password", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		errors = addError("password", "min:6", errors, tmpError)
	}

	if match, _ := regexp.MatchString("[0-9]", value); !match {
		tmpError := "La contraseña debe contener al menos un número"
		tmpErrorKey := fmt.Sprintf("%s.password", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		errors = addError("password", "regex", errors, tmpError)
	}

	if match, _ := regexp.MatchString("[a-z]", value); !match {
		tmpError := "La contraseña debe contener al menos una letra minúscula"
		tmpErrorKey := fmt.Sprintf("%s.password", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		errors = addError("password", "regex", errors, tmpError)
	}

	if match, _ := regexp.MatchString("[A-Z]", value); !match {
		tmpError := "La contraseña debe contener al menos una letra mayúscula"
		tmpErrorKey := fmt.Sprintf("%s.password", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		errors = addError("password", "regex", errors, tmpError)
	}

	if match, _ := regexp.MatchString(`[#$%&/()!_-]+`, value); !match {
		tmpError := "La contraseña debe contener al menos un carácter especial ($#%&/()!_-)"
		tmpErrorKey := fmt.Sprintf("%s.password", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		errors = addError("password", "regex", errors, tmpError)
	}

	return errors
}
