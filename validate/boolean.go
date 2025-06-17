package validate

import (
	"log"
	"strconv"
)

func Boolean(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}) map[string]interface{} {
	if _, ok := payload[input]; !ok {
		return errors
	}

	value := payload[input]
	log.Printf("Validating boolean for input '%s' with value '%v'\n", input, value)

	_, err := strconv.ParseBool(value.(string))
	if err != nil {
		log.Printf("Error parsing boolean for input '%s': %v\n", input, err)
		errors = addError(input, "boolean", errors, "El valor debe ser booleano (true/false).")
	}

	return errors
}
