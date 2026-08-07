package validate

// addError de prueba: replica el comportamiento del addError real de govalidator
// (paquete distinto, no exportado, por lo que no se puede reutilizar directamente).
func testAddError(input string, rule string, allErrors map[string]interface{}, err string) map[string]interface{} {
	if _, exists := allErrors[input]; !exists {
		allErrors[input] = map[string]interface{}{
			rule: []string{err},
		}
		return allErrors
	}

	fieldErrors, ok := allErrors[input].(map[string]interface{})
	if !ok {
		allErrors[input] = map[string]interface{}{
			rule: []string{err},
		}
		return allErrors
	}

	if _, exists := fieldErrors[rule]; !exists {
		fieldErrors[rule] = []string{err}
	} else {
		fieldErrors[rule] = append(fieldErrors[rule].([]string), err)
	}
	allErrors[input] = fieldErrors

	return allErrors
}

func getErrorMsgs(errors map[string]interface{}, input, rule string) ([]string, bool) {
	fieldErrors, ok := errors[input].(map[string]interface{})
	if !ok {
		return nil, false
	}

	msgs, ok := fieldErrors[rule].([]string)
	return msgs, ok
}
