package validate

import (
	"strconv"
	"time"
)

func After(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}) map[string]interface{} {
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
			if date.Before(time.Now()) {
				errors = addError(input, "after", errors, "La fecha no es posterior a la fecha actual")
			}
		}

		if options[0] == "tomorrow" {
			if date.Before(time.Now().AddDate(0, 0, 1)) {
				errors = addError(input, "after", errors, "La fecha no es posterior a mañana")
			}
		}

		if options[0] == "yesterday" {
			if date.Before(time.Now().AddDate(0, 0, -1)) {
				errors = addError(input, "after", errors, "La fecha no es posterior a ayer")
			}
		}

		compare_date, err := time.Parse("2006-01-02", options[0])
		if err == nil {
			if date.Before(compare_date) {
				errors = addError(input, "after", errors, "La fecha no es posterior a la fecha "+options[0])
			}
		}
	}

	num, ok := value.(int)
	if ok {
		if len(options) == 0 {
			errors = addError(input, "after", errors, "El valor de comparación no está definido")
		} else {
			compare_int, _ := strconv.Atoi(options[0])
			if num < compare_int {
				errors = addError(input, "after", errors, "El valor de "+input+" no es mayor que "+options[0])
			}
		}
	}

	return errors
}
