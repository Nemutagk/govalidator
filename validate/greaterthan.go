package validate

import (
	"fmt"
)

func GreaterThan(input string, value any, payload map[string]any, options []string, sliceIndex string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, customeErrors map[string]string) map[string]interface{} {
	if len(options) == 0 {
		errors = addError(input, "greater_than", errors, "El valor a comparar no está definido")
		return errors
	}

	pair, errors, ok := resolveComparable(input, value, payload, options, "greater_than", sliceIndex, errors, addError, customeErrors)
	if !ok {
		return errors
	}

	var failed bool
	if pair.isDate {
		failed = !pair.ownDate.After(pair.cmpDate)
	} else {
		failed = pair.ownNum <= pair.cmpNum
	}

	if failed {
		tmpError := fmt.Sprintf("El campo %s debe ser mayor que %s", input, options[0])

		if sliceIndex != "" {
			tmpError = fmt.Sprintf("El campo %s en la posición %s debe ser mayor que %s", input, sliceIndex, options[0])
		}

		tmpErrorKey := fmt.Sprintf("%s.greater_than", input)
		if customeError, exists := customeErrors[tmpErrorKey]; exists {
			tmpError = customeError
		}
		errors = addError(input, "greater_than", errors, tmpError)
	}

	return errors
}
