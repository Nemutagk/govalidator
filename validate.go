package govalidator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Nemutagk/govalidator/db"
	"github.com/Nemutagk/govalidator/validate"
)

func ValidateRequestPrepare(r http.Request) (map[string]interface{}, error) {
	var payload map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("error to decode payload: %v", err)
	}

	return payload, nil
}

func ValidateRequest(body map[string]interface{}, rules map[string]string, dbManager *db.ConnectionManager) (map[string]interface{}, map[string]interface{}) {

	rules_intpus := make(map[string]interface{})

	body_encode, err := json.Marshal(body)

	if err != nil {
		return nil, map[string]interface{}{
			"error": "The payload in invalid!",
		}
	}

	var body_parse map[string]interface{}

	err = json.Unmarshal(body_encode, &body_parse)

	if err != nil {
		return nil, map[string]interface{}{
			"error": "The payload in invalid!",
		}
	}

	for input, rules := range rules {
		separated_rules := strings.Split(rules, "|")
		rules_building := make(map[string]interface{})
		for _, rule := range separated_rules {
			tmp_rule := getRuleWithOptions(rule)
			rules_building[tmp_rule[0]] = tmp_rule[1:]
		}
		rules_intpus[input] = rules_building
	}

	errors := make(map[string]interface{})
	safePayload := make(map[string]interface{})
	for input, rules := range rules_intpus {
		rulesMap, ok := rules.(map[string]interface{})
		if !ok {
			continue
		}
		for rule, options := range rulesMap {
			opts, ok := options.([]string)
			if !ok {
				continue
			}
			switch rule {
			case "required":
				errors = validate.Required(input, body_parse, opts, errors, addError)
			case "email":
				errors = validate.Email(input, body_parse, opts, errors, addError)
			case "confirmation":
				errors = validate.Confirmation(input, body_parse, opts, errors, addError)
			case "unique":
				errors = validate.Unique(input, body_parse, opts, errors, addError, dbManager)
			case "in":
				errors = validate.In(input, body_parse, opts, errors, addError)
			case "before":
				errors = validate.Before(input, body_parse, opts, errors, addError)
			case "after":
				errors = validate.After(input, body_parse, opts, errors, addError)
			case "ip":
				errors = validate.Ip(input, body_parse, opts, errors, addError)
			case "password":
				errors = validate.Password(input, body_parse, opts, errors, addError)
			case "exists":
				errors = validate.Exists(input, body_parse, opts, errors, addError, dbManager)
			case "min":
				errors = validate.Min(input, body_parse, opts, errors, addError)
			case "max":
				errors = validate.Max(input, body_parse, opts, errors, addError)
			case "boolean":
				errors = validate.Boolean(input, body_parse, opts, errors, addError)
			case "sometimes":
				fmt.Println("Sometimes rule is not implemented yet")
			case "required_with":
				errors = validate.RequiredWith(input, body_parse, opts, errors, addError)
			case "required_with_all":
				errors = validate.RequiredWithAll(input, body_parse, opts, errors, addError)
			case "required_without":
				errors = validate.RequiredWithout(input, body_parse, opts, errors, addError)
			case "required_without_all":
				errors = validate.RequiredWithoutAll(input, body_parse, opts, errors, addError)

			default:
				errors = addError(input, input, errors, "The rule "+rule+" is not valid")
			}
		}

		if _, ok := body_parse[input]; ok {
			safePayload[input] = body_parse[input]
		}
	}

	return safePayload, errors
}

func getRuleWithOptions(rule string) []string {
	parts := strings.Split(rule, ":")

	var build_rule []string
	build_rule = append(build_rule, parts[0])

	if len(parts) > 1 {
		options := strings.Split(parts[1], ",")

		build_rule = append(build_rule, options...)
	}

	return build_rule
}

func addError(input string, rule string, errors map[string]interface{}, error string) map[string]interface{} {
	if _, exists_input := errors[input]; !exists_input {
		errors[input] = map[string]interface{}{
			rule: []string{
				error,
			},
		}
	} else {
		if inputErrors, ok := errors[input].(map[string]interface{}); ok {
			if _, exists_rule := inputErrors[rule]; !exists_rule {
				inputErrors[rule] = []string{
					error,
				}
			} else {
				inputErrors[rule] = append(inputErrors[rule].([]string), error)
			}
			errors[input] = inputErrors
		} else {
			errors[input] = map[string]interface{}{
				rule: []string{
					error,
				},
			}
		}
	}

	return errors
}
