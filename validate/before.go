package validate

import (
	"strconv"
	"time"
)

func Before(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}) map[string]interface{} {
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
			if date.After(time.Now()) {
				errors = addError(input, "before", errors, "The date is not before the current date")
			}
		}

		if options[0] == "tomorrow" {
			if date.After(time.Now().AddDate(0, 0, 1)) {
				errors = addError(input, "before", errors, "The date is not before tomorrow")
			}
		}

		if options[0] == "yesterday" {
			if date.After(time.Now().AddDate(0, 0, -1)) {
				errors = addError(input, "before", errors, "The date is not before yesterday")
			}
		}

		compare_date, err := time.Parse("2006-01-02", options[0])
		if err == nil {
			if date.After(compare_date) {
				errors = addError(input, "before", errors, "The date is not before the date "+options[0])
			}
		}
	}

	num, ok := value.(int)
	if ok {
		if len(options) == 0 {
			errors = addError(input, "before", errors, "The compare value is not defined")
		} else {
			compare_int, _ := strconv.Atoi(options[0])
			if num > compare_int {
				errors = addError(input, "before", errors, "The input "+input+" is not before the number "+options[0])
			}
		}
	}

	return errors
}
