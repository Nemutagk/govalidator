package validate

import "strconv"

func Boolean(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}) map[string]interface{} {
	if _, ok := payload[input]; !ok {
		return errors
	}

	_, err := strconv.ParseBool(payload[input].(string))
	if err != nil {
		errors = addError(input, "boolean", errors, "The input "+input+" is not a valid boolean: "+err.Error())
	}

	return errors
}
