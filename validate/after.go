package validate

import (
	"strconv"
	"time"
)

func After(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}) map[string]interface{} {
	if len(options) == 0 {
		errors = addError(input, "before", errors, "The compare value is not defined")
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
				errors = addError(input, "after", errors, "The date is not after the current date")
			}
		}

		if options[0] == "tomorrow" {
			if date.Before(time.Now().AddDate(0, 0, 1)) {
				errors = addError(input, "after", errors, "The date is not after tomorrow")
			}
		}

		if options[0] == "yesterday" {
			if date.Before(time.Now().AddDate(0, 0, -1)) {
				errors = addError(input, "after", errors, "The date is not after yesterday")
			}
		}

		compare_date, err := time.Parse("2006-01-02", options[0])
		if err == nil {
			if date.Before(compare_date) {
				errors = addError(input, "after", errors, "The date is not after the date "+options[0])
			}
		}
	}

	num, ok := value.(int)
	if ok {
		if len(options) == 0 {
			errors = addError(input, "after", errors, "The compare value is not defined")
		} else {
			compare_int, _ := strconv.Atoi(options[0])
			if num < compare_int {
				errors = addError(input, "after", errors, "The input "+input+" is not after the number "+options[0])
			}
		}
	}

	return errors
}
