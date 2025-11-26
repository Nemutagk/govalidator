package govalidator

import (
	"strconv"
	"strings"

	"github.com/Nemutagk/goerrors"
	"github.com/Nemutagk/govalidator/v2/validate"
)

type Input struct {
	Name   string
	Parent string
	Rules  []Rule
}

type Rule struct {
	Name    string
	Options []string
}

func ValidateRequest(body map[string]any, inputs []Input, customeallErrors map[string]string, models map[string]func(data string, payload map[string]any) bool) (map[string]any, *goerrors.GError) {
	safePayload, currentallErrors, _ := rangeInputs(body, inputs, customeallErrors, models, "")

	allErrors := make([]string, 0)
	if len(currentallErrors) > 0 {
		for input, inputallErrors := range currentallErrors {
			for _, errMessages := range inputallErrors.(map[string]interface{}) {
				for _, errMessage := range errMessages.([]string) {
					allErrors = append(allErrors, input+": "+errMessage)
				}
			}
		}
		gerr := goerrors.NewGError("Error en la validación", goerrors.StatusBadRequest, &allErrors, nil)

		return nil, gerr
	}

	return safePayload, nil
}

func rangeInputs(body map[string]any, inputs []Input, customeallErrors map[string]string, models map[string]func(data string, payload map[string]any) bool, sliceIndex string) (map[string]any, map[string]any, map[string]bool) {
	// log.Printf("====================> Starting rangeInputs ====================")
	safePayload := make(map[string]any)
	allErrors := make(map[string]any)
	includesSometimesRule := make(map[string]bool)

	for index, input := range inputs {
		// log.Printf("--------------------------------------------------")
		inputName := input.Name
		// log.Printf("Processing raw input: %s", inputName)

		if strings.Contains(inputName, ".") {
			parts := strings.Split(inputName, ".")
			inputName = parts[0]
			input.Name = strings.Join(parts[1:], ".")

			if input.Parent == "" {
				input.Parent = inputName
			}
		}
		// log.Printf("Processed input: %s", inputName)

		value, exists := body[inputName]
		if !exists {
			value = nil
		}

		if input.Parent != "" {
			if _, exists := includesSometimesRule[input.Parent]; exists {
				input.Rules = append([]Rule{{Name: "sometimes"}}, input.Rules...)
			}
		}

		switch value.(type) {
		case map[string]any:
			tmpPayload, tmpErrors, tmpSometimes := rangeInputs(value.(map[string]any), []Input{input}, customeallErrors, models, sliceIndex)
			for k, v := range tmpPayload {
				safePayload[k] = v
			}
			for k, v := range tmpErrors {
				allErrors[k] = v
			}
			for k, v := range tmpSometimes {
				includesSometimesRule[k] = v
			}
		case []any:
			// log.Printf("Input %s is an array", input.Name)
			// log.Printf("inputName: %s", inputName)
			sliceIndexStr := strconv.Itoa(index)
			if inputName != input.Name {
				tmpPayload, tmpErrors, tmpSometimes := rangeArrayInput(value.([]any), body, input, customeallErrors, models, sliceIndexStr)
				if safePayload[inputName] == nil {
					safePayload[inputName] = []any{}
				}

				valueSlice, err := value.([]any)
				if !err {
					valueSlice = []any{}
				}

				for k, v := range tmpPayload {
					if k < len(valueSlice) {
						valueSlice[k] = v
					} else {
						valueSlice = append(valueSlice, v)
					}
				}

				for k, v := range tmpErrors {
					allErrors[k] = v
				}
				for k, v := range tmpSometimes {
					includesSometimesRule[k] = v
				}
				continue
			}

			tmpPayload, tmpErrors, tmpSometimes := applyRules(inputName, input, value, body, customeallErrors, models, sliceIndexStr)
			safePayload[input.Name] = tmpPayload
			for k, v := range tmpErrors {
				allErrors[k] = v
			}
			for k, v := range tmpSometimes {
				includesSometimesRule[k] = v
			}
		default:
			value, tmpErrors, tmpSometimes := applyRules(input.Name, input, value, body, customeallErrors, models, sliceIndex)
			safePayload[input.Name] = value
			for k, v := range tmpErrors {
				allErrors[k] = v
			}
			for k, v := range tmpSometimes {
				includesSometimesRule[k] = v
			}
		}
	}

	return safePayload, allErrors, includesSometimesRule
}

func rangeArrayInput(body []any, fullBody map[string]any, inputs Input, customeallErrors map[string]string, models map[string]func(data string, payload map[string]any) bool, sliceIndex string) ([]any, map[string]any, map[string]bool) {
	// log.Printf("====================> Starting rangeArrayInput ====================")
	safePayload := []any{}
	allErrors := make(map[string]any)
	includesSometimesRule := make(map[string]bool)

	inputName := inputs.Name
	if strings.Contains(inputName, ".") {
		parts := strings.Split(inputName, ".")
		inputName = parts[0]
		inputs.Name = strings.Join(parts[1:], ".")

		if inputs.Parent == "" {
			inputs.Parent = inputName
		}
	}
	// log.Printf("Processing array input: %s", inputName)

	for index, value := range body {
		if inputName == "*" {
			indexStr := strconv.Itoa(index)
			value, tmpError, tmpSometimes := applyRules(inputs.Name, inputs, value, value.(map[string]any), customeallErrors, models, indexStr)
			safePayload = append(safePayload, value)
			for k, v := range tmpError {
				allErrors[k] = v
			}
			for k, v := range tmpSometimes {
				includesSometimesRule[k] = v
			}
		} else {
			indexInt, err := strconv.Atoi(inputName)
			if err != nil {
				// log.Printf("Error converting index to int: %v", err)
				continue
			}

			if indexInt == index {
				indexStr := strconv.Itoa(index)
				value, tmpError, tmpSometimes := applyRules(indexStr, inputs, value, value.(map[string]any), customeallErrors, models, indexStr)
				safePayload = append(safePayload, value)
				for k, v := range tmpError {
					allErrors[k] = v
				}
				for k, v := range tmpSometimes {
					includesSometimesRule[k] = v
				}
			}
		}
	}

	return safePayload, allErrors, includesSometimesRule
}

func applyRules(inputName any, input Input, value any, body map[string]any, customeallErrors map[string]string, models map[string]func(data string, payload map[string]any) bool, sliceIndex string) (any, map[string]any, map[string]bool) {
	// log.Printf("====================> Starting applyRules for input ====================")
	// log.Printf("Input Name: %v", inputName)
	// log.Printf("Input Value: %+v", value)
	// log.Printf("Input Rules: %+v", input.Rules)
	// log.Printf("sliceIndex : %s", sliceIndex)

	allErrors := make(map[string]any)
	includesSometimesRule := make(map[string]bool)

	skipRulesMap := false

	inputNameStr := ""
	if _, ok := inputName.(string); !ok {
		inputNameStr = strconv.Itoa(inputName.(int))
	} else {
		inputNameStr = inputName.(string)
	}

	for _, rule := range input.Rules {
		opts := rule.Options

		switch rule.Name {
		case "email":
			allErrors = validate.Email(inputNameStr, value, body, opts, sliceIndex, allErrors, addError, customeallErrors)
		case "confirmation":
			allErrors = validate.Confirmation(inputNameStr, value, body, opts, sliceIndex, allErrors, addError, customeallErrors)
		case "unique":
			allErrors = validate.Unique(inputNameStr, value, body, opts, sliceIndex, allErrors, addError, models, customeallErrors)
		case "in":
			allErrors = validate.In(inputNameStr, value, body, opts, sliceIndex, allErrors, addError, customeallErrors)
		case "not_in":
			allErrors = validate.NotIn(inputNameStr, value, body, opts, sliceIndex, allErrors, addError, customeallErrors)
		case "before":
			allErrors = validate.Before(inputNameStr, value, body, opts, sliceIndex, allErrors, addError, customeallErrors)
		case "after":
			allErrors = validate.After(inputNameStr, value, body, opts, sliceIndex, allErrors, addError, customeallErrors)
		case "ip":
			allErrors = validate.Ip(inputNameStr, value, body, opts, sliceIndex, allErrors, addError, customeallErrors)
		case "password":
			allErrors = validate.Password(inputNameStr, value, body, opts, sliceIndex, allErrors, addError, customeallErrors)
		case "exists":
			allErrors = validate.Exists(inputNameStr, value, body, opts, sliceIndex, allErrors, addError, models, customeallErrors)
		case "min":
			allErrors = validate.Min(inputNameStr, value, body, opts, sliceIndex, allErrors, addError, customeallErrors)
		case "max":
			allErrors = validate.Max(inputNameStr, value, body, opts, sliceIndex, allErrors, addError, customeallErrors)
		case "boolean":
			allErrors = validate.Boolean(inputNameStr, value, body, opts, sliceIndex, allErrors, addError, customeallErrors)
		case "sometimes":
			if _, exists_input := body[inputNameStr]; !exists_input {
				//Si no existe el input no se validan lás demás reglas existentes
				// // // log.Println("Input", input, "no existe en el body, no se validan las demás reglas")
				skipRulesMap = true
				includesSometimesRule[inputNameStr] = true
			}
		case "required":
			allErrors = validate.Required(inputNameStr, value, body, opts, sliceIndex, allErrors, addError, customeallErrors)
		case "required_with":
			allErrors = validate.RequiredWith(inputNameStr, value, body, opts, sliceIndex, allErrors, addError, customeallErrors)
		case "required_with_all":
			allErrors = validate.RequiredWithAll(inputNameStr, value, body, opts, sliceIndex, allErrors, addError, customeallErrors)
		case "required_without":
			allErrors = validate.RequiredWithout(inputNameStr, value, body, opts, sliceIndex, allErrors, addError, customeallErrors)
		case "required_without_all":
			allErrors = validate.RequiredWithoutAll(inputNameStr, value, body, opts, sliceIndex, allErrors, addError, customeallErrors)
		case "array":
			allErrors = validate.Array(inputNameStr, value, body, opts, sliceIndex, allErrors, addError, customeallErrors)
		case "type":
			allErrors = validate.Type(inputNameStr, value, body, opts, sliceIndex, allErrors, addError, customeallErrors)
		case "date":
			allErrors = validate.Date(inputNameStr, value, body, opts, sliceIndex, allErrors, addError, customeallErrors)
		case "date_format":
			allErrors = validate.DateFormat(inputNameStr, value, body, opts, sliceIndex, allErrors, addError, customeallErrors)
		case "custome":
			allErrors = validate.Custome(inputNameStr, value, body, opts, sliceIndex, allErrors, addError, models, customeallErrors)

		default:
			allErrors = addError(inputNameStr, rule.Name, allErrors, "The rule "+rule.Name+" is not valid")
		}

		if skipRulesMap {
			break
		}
	}

	return value, allErrors, includesSometimesRule
}

func addError(input string, rule string, allErrors map[string]interface{}, error string) map[string]interface{} {
	if _, exists_input := allErrors[input]; !exists_input {
		allErrors[input] = map[string]interface{}{
			rule: []string{
				error,
			},
		}
	} else {
		if inputallErrors, ok := allErrors[input].(map[string]interface{}); ok {
			if _, exists_rule := inputallErrors[rule]; !exists_rule {
				inputallErrors[rule] = []string{
					error,
				}
			} else {
				inputallErrors[rule] = append(inputallErrors[rule].([]string), error)
			}
			allErrors[input] = inputallErrors
		} else {
			allErrors[input] = map[string]interface{}{
				rule: []string{
					error,
				},
			}
		}
	}

	return allErrors
}
