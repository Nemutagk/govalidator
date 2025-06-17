package validate

import (
	"strconv"
	"time"
)

func Before(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}) map[string]interface{} {
	if len(options) == 0 {
		errors = addError(input, "before", errors, "El valor de comparación no está definido")
		return errors
	}

	value := payload[input]

	if value == nil || value == "" {
		return errors
	}

	date, err_date := time.Parse("2006-01-02", value.(string))

	if err_date == nil {
		if options[0] == "now" || options[0] == "current" || options[0] == "today" {
			if date.After(time.Now()) {
				errors = addError(input, "before", errors, "La fecha no es anterior a la fecha actual")
			}
		}

		if options[0] == "tomorrow" {
			if date.After(time.Now().AddDate(0, 0, 1)) {
				errors = addError(input, "before", errors, "La fecha no es anterior a mañana")
			}
		}

		if options[0] == "yesterday" {
			if date.After(time.Now().AddDate(0, 0, -1)) {
				errors = addError(input, "before", errors, "La fecha no es anterior a ayer")
			}
		}

		compare_date, err := time.Parse("2006-01-02", options[0])
		if err == nil {
			if date.After(compare_date) {
				errors = addError(input, "before", errors, "La fecha no es anterior a la fecha "+options[0])
			}
		}
	}

	num, ok := value.(int)
	if ok {
		if len(options) == 0 {
			errors = addError(input, "before", errors, "El valor de comparación no está definido")
		} else {
			compare_int, _ := strconv.Atoi(options[0])
			if num > compare_int {
				errors = addError(input, "before", errors, "El valor de "+input+" no es menor que "+options[0])
			}
		}
	}

	return errors
}
