package validate

import (
	"fmt"
	"regexp"
)

func Password(input string, value any, payload map[string]any, options []string, sliceIndex string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, customeErrors map[string]string) map[string]interface{} {
	if _, exists_input := payload[input]; !exists_input {
		return errors
	}

	valueStr, ok := value.(string)
	if !ok {
		tmpError := "La contraseña debe ser una cadena de texto"

		if sliceIndex != "" {
			tmpError = fmt.Sprintf("La contraseña en la posición %s debe ser una cadena de texto", sliceIndex)
		}

		tmpErrorKey := fmt.Sprintf("%v.password", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		errors = addError("password", "type", errors, tmpError)
		return errors
	}

	if len(valueStr) < 6 {
		tmpError := "La contraseña debe tener al menos 6 caracteres"

		if sliceIndex != "" {
			tmpError = fmt.Sprintf("La contraseña en la posición %s debe tener al menos 6 caracteres", sliceIndex)
		}

		tmpErrorKey := fmt.Sprintf("%s.password", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		errors = addError("password", "min:6", errors, tmpError)
	}

	if match, _ := regexp.MatchString("[0-9]", valueStr); !match {
		tmpError := "La contraseña debe contener al menos un número"

		if sliceIndex != "" {
			tmpError = fmt.Sprintf("La contraseña en la posición %s debe contener al menos un número", sliceIndex)
		}

		tmpErrorKey := fmt.Sprintf("%s.password", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		errors = addError("password", "regex", errors, tmpError)
	}

	if match, _ := regexp.MatchString("[a-z]", valueStr); !match {
		tmpError := "La contraseña debe contener al menos una letra minúscula"

		if sliceIndex != "" {
			tmpError = fmt.Sprintf("La contraseña en la posición %s debe contener al menos una letra minúscula", sliceIndex)
		}

		tmpErrorKey := fmt.Sprintf("%s.password", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		errors = addError("password", "regex", errors, tmpError)
	}

	if match, _ := regexp.MatchString("[A-Z]", valueStr); !match {
		tmpError := "La contraseña debe contener al menos una letra mayúscula"

		if sliceIndex != "" {
			tmpError = fmt.Sprintf("La contraseña en la posición %s debe contener al menos una letra mayúscula", sliceIndex)
		}

		tmpErrorKey := fmt.Sprintf("%s.password", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		errors = addError("password", "regex", errors, tmpError)
	}

	if match, _ := regexp.MatchString(`[#$%&/()!_-]+`, valueStr); !match {
		tmpError := "La contraseña debe contener al menos un carácter especial ($#%&/()!_-)"

		if sliceIndex != "" {
			tmpError = fmt.Sprintf("La contraseña en la posición %s debe contener al menos un carácter especial ($#%&/()!_-)", sliceIndex)
		}

		tmpErrorKey := fmt.Sprintf("%s.password", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		errors = addError("password", "regex", errors, tmpError)
	}

	return errors
}
