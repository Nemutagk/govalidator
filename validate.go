package govalidator

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

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

type ValidationError struct {
	FieldErrors []string
}

func (u ValidationError) Error() string {
	return "Se encontraron errores en la validación"
}

func ValidateRequest(body map[string]any, inputs []Input, customeallErrors map[string]string, models map[string]func(data any, payload map[string]any, opts *[]string) (bool, string)) (map[string]any, error) {
	safePayload, currentallErrors, _ := rangeInputs(body, inputs, customeallErrors, models, "", "", body)

	allErrors := make([]string, 0)
	if len(currentallErrors) > 0 {
		for input, inputallErrors := range currentallErrors {
			for _, errMessages := range inputallErrors.(map[string]interface{}) {
				for _, errMessage := range errMessages.([]string) {
					allErrors = append(allErrors, fmt.Sprintf("%s: %s", input, errMessage))
				}
			}
		}

		return nil, ValidationError{FieldErrors: allErrors}
	}

	for k := range safePayload {
		if k == "" {
			// log.Printf("el campo tiene nombre vacío, se elimina del payload seguro")
			delete(safePayload, k)
			continue
		}

		val, exists := body[k]
		if !exists {
			// log.Printf("el campo %s no existe en el payload original", k)
			delete(safePayload, k)
			continue
		} else {
			if val == nil && (!existsRule(inputs, k, "sometimes") && !existsRule(inputs, k, "nullable")) {
				// log.Printf("el campo %s es nil y no tiene la regla sometimes o nullable", k)
				delete(safePayload, k)
				continue
			}
		}
		// log.Printf("el campo %s existe en el payload original", k)
	}

	return safePayload, nil
}

func ConvertPayload[T any](payload map[string]any) (T, error) {
	var result T
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(payloadBytes, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func rangeInputs(body map[string]any, inputs []Input, customeallErrors map[string]string, models map[string]func(data any, payload map[string]any, opts *[]string) (bool, string), sliceIndex string, pathPrefix string, rootBody map[string]any) (map[string]any, map[string]any, map[string]bool) {
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

		switch value := value.(type) {
		case map[string]any:
			newPrefix := inputName
			if pathPrefix != "" {
				newPrefix = pathPrefix + "." + inputName
			}
			tmpPayload, tmpErrors, tmpSometimes := rangeInputs(value, []Input{input}, customeallErrors, models, sliceIndex, newPrefix, rootBody)
			if _, ok := safePayload[inputName]; !ok {
				safePayload[inputName] = make(map[string]any)
			}
			deepMerge(safePayload[inputName].(map[string]any), tmpPayload)
			for k, v := range tmpErrors {
				allErrors[k] = v
			}
			for k, v := range tmpSometimes {
				includesSometimesRule[k] = v
			}
		case []any:
			sliceIndexStr := strconv.Itoa(index)
			if inputName != input.Name {
				remaining := input.Name
				specificIndex := -1

				if strings.HasPrefix(remaining, "*.") {
					remaining = remaining[2:]
				} else if remaining == "*" {
					remaining = ""
				} else {
					parts := strings.SplitN(remaining, ".", 2)
					idx, err := strconv.Atoi(parts[0])
					if err != nil {
						continue
					}
					specificIndex = idx
					if len(parts) == 2 {
						remaining = parts[1]
					} else {
						remaining = ""
					}
				}

				arrayPrefix := inputName
				if pathPrefix != "" {
					arrayPrefix = pathPrefix + "." + inputName
				}

				resultSlice := make([]any, 0, len(value))
				for i, elem := range value {
					if specificIndex >= 0 && i != specificIndex {
						continue
					}

					elemPrefix := arrayPrefix + "." + strconv.Itoa(i)

					if remaining == "" {
						indexStr := strconv.Itoa(i)
						syntheticBody := map[string]any{indexStr: elem}
						_, tmpErrors, tmpSometimes := applyRules(indexStr, input, elem, syntheticBody, customeallErrors, models, indexStr, rootBody)
						for k, v := range tmpErrors {
							allErrors[arrayPrefix+"."+k] = v
						}
						for k, v := range tmpSometimes {
							includesSometimesRule[arrayPrefix+"."+k] = v
						}
						resultSlice = append(resultSlice, elem)
					} else {
						var elemBody map[string]any
						if m, ok := elem.(map[string]any); ok {
							elemBody = m
						} else {
							elemBody = make(map[string]any)
						}
						elemInput := Input{Name: remaining, Rules: input.Rules, Parent: input.Parent}
						tmpPayload, tmpErrors, tmpSometimes := rangeInputs(elemBody, []Input{elemInput}, customeallErrors, models, strconv.Itoa(i), elemPrefix, rootBody)
						resultSlice = append(resultSlice, tmpPayload)
						for k, v := range tmpErrors {
							allErrors[k] = v
						}
						for k, v := range tmpSometimes {
							includesSometimesRule[k] = v
						}
					}
				}

				mergeOffset := 0
				if specificIndex >= 0 {
					mergeOffset = specificIndex
				}
				if existing, ok := safePayload[inputName].([]any); ok {
					for i, elem := range resultSlice {
						targetI := i + mergeOffset
						if targetI < len(existing) {
							if existingMap, ok := existing[targetI].(map[string]any); ok {
								if elemMap, ok := elem.(map[string]any); ok {
									deepMerge(existingMap, elemMap)
									continue
								}
							}
							existing[targetI] = elem
						} else {
							existing = append(existing, elem)
						}
					}
					safePayload[inputName] = existing
				} else {
					safePayload[inputName] = resultSlice
				}
				continue
			}

			tmpPayload, tmpErrors, tmpSometimes := applyRules(inputName, input, value, body, customeallErrors, models, sliceIndexStr, rootBody)
			safePayload[input.Name] = tmpPayload
			for k, v := range tmpErrors {
				if pathPrefix != "" {
					allErrors[pathPrefix+"."+k] = v
				} else {
					allErrors[k] = v
				}
			}
			for k, v := range tmpSometimes {
				if pathPrefix != "" {
					includesSometimesRule[pathPrefix+"."+k] = v
				} else {
					includesSometimesRule[k] = v
				}
			}
		default:
			if inputName != input.Name {
				if strings.HasPrefix(input.Name, "*.") || input.Name == "*" {
					break
				}
				newPrefix := inputName
				if pathPrefix != "" {
					newPrefix = pathPrefix + "." + inputName
				}
				tmpPayload, tmpErrors, tmpSometimes := rangeInputs(make(map[string]any), []Input{input}, customeallErrors, models, sliceIndex, newPrefix, rootBody)
				if _, ok := safePayload[inputName]; !ok {
					safePayload[inputName] = make(map[string]any)
				}
				if existing, ok := safePayload[inputName].(map[string]any); ok {
					deepMerge(existing, tmpPayload)
				}
				for k, v := range tmpErrors {
					allErrors[k] = v
				}
				for k, v := range tmpSometimes {
					includesSometimesRule[k] = v
				}
			} else {
				value, tmpErrors, tmpSometimes := applyRules(input.Name, input, value, body, customeallErrors, models, sliceIndex, rootBody)
				safePayload[input.Name] = value
				for k, v := range tmpErrors {
					if pathPrefix != "" {
						allErrors[pathPrefix+"."+k] = v
					} else {
						allErrors[k] = v
					}
				}
				for k, v := range tmpSometimes {
					if pathPrefix != "" {
						includesSometimesRule[pathPrefix+"."+k] = v
					} else {
						includesSometimesRule[k] = v
					}
				}
			}
		}
	}

	return safePayload, allErrors, includesSometimesRule
}

func applyRules(inputName any, input Input, value any, body map[string]any, customeallErrors map[string]string, models map[string]func(data any, payload map[string]any, opts *[]string) (bool, string), sliceIndex string, rootBody map[string]any) (any, map[string]any, map[string]bool) {
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
			getError := false
			allErrors, getError = validate.Required(inputNameStr, value, body, opts, sliceIndex, allErrors, addError, customeallErrors)

			if getError {
				skipRulesMap = true
				includesSometimesRule[inputNameStr] = true
			}
		case "required_with":
			getError := false
			allErrors, getError = validate.RequiredWith(inputNameStr, value, body, opts, sliceIndex, allErrors, addError, customeallErrors)
			if getError {
				skipRulesMap = true
				includesSometimesRule[inputNameStr] = true
			}
		case "required_with_all":
			getError := false
			allErrors, getError = validate.RequiredWithAll(inputNameStr, value, body, opts, sliceIndex, allErrors, addError, customeallErrors)
			if getError {
				skipRulesMap = true
				includesSometimesRule[inputNameStr] = true
			}
		case "required_without":
			getError := false
			allErrors, getError = validate.RequiredWithout(inputNameStr, value, body, opts, sliceIndex, allErrors, addError, customeallErrors)
			if getError {
				skipRulesMap = true
				includesSometimesRule[inputNameStr] = true
			}
		case "required_without_all":
			getError := false
			allErrors, getError = validate.RequiredWithoutAll(inputNameStr, value, body, opts, sliceIndex, allErrors, addError, customeallErrors)
			if getError {
				skipRulesMap = true
				includesSometimesRule[inputNameStr] = true
			}
		case "array":
			allErrors = validate.Array(inputNameStr, value, body, opts, sliceIndex, allErrors, addError, customeallErrors)
		case "type":
			allErrors = validate.Type(inputNameStr, value, body, opts, sliceIndex, allErrors, addError, customeallErrors)
		case "date":
			allErrors = validate.Date(inputNameStr, value, body, opts, sliceIndex, allErrors, addError, customeallErrors)
		case "date_format":
			allErrors = validate.DateFormat(inputNameStr, value, body, opts, sliceIndex, allErrors, addError, customeallErrors)
		case "customized":
			allErrors = validate.Customized(inputNameStr, value, body, opts, sliceIndex, allErrors, addError, models, customeallErrors)
		case "nullable":
			allErrors = validate.Nullable(inputNameStr, value, body, opts, sliceIndex, allErrors, addError, customeallErrors)
		case "equal":
			allErrors = validate.Equal(inputNameStr, value, body, opts, sliceIndex, allErrors, addError, customeallErrors)
		case "not_equal":
			allErrors = validate.NotEqual(inputNameStr, value, body, opts, sliceIndex, allErrors, addError, customeallErrors)
		case "required_if":
			allErrors = validate.RequiredIf(inputNameStr, value, rootBody, opts, sliceIndex, allErrors, addError, customeallErrors)
		default:
			allErrors = addError(inputNameStr, rule.Name, allErrors, "The rule "+rule.Name+" is not valid")
		}

		if skipRulesMap {
			break
		}
	}

	return value, allErrors, includesSometimesRule
}

func addError(input string, rule string, allErrors map[string]interface{}, err string) map[string]interface{} {
	if _, exists_input := allErrors[input]; !exists_input {
		allErrors[input] = map[string]interface{}{
			rule: []string{
				err,
			},
		}
	} else {
		if inputallErrors, ok := allErrors[input].(map[string]interface{}); ok {
			if _, exists_rule := inputallErrors[rule]; !exists_rule {
				inputallErrors[rule] = []string{
					err,
				}
			} else {
				inputallErrors[rule] = append(inputallErrors[rule].([]string), err)
			}
			allErrors[input] = inputallErrors
		} else {
			allErrors[input] = map[string]interface{}{
				rule: []string{
					err,
				},
			}
		}
	}

	return allErrors
}

func deepMerge(dst, src map[string]any) {
	for k, v := range src {
		if existing, ok := dst[k]; ok {
			if existingMap, ok := existing.(map[string]any); ok {
				if srcMap, ok := v.(map[string]any); ok {
					deepMerge(existingMap, srcMap)
					continue
				}
			}
		}
		dst[k] = v
	}
}

func existsRule(inputs []Input, inputName, rule string) bool {
	for _, inp := range inputs {
		if inp.Name == inputName {
			for _, r := range inp.Rules {
				if r.Name == rule {
					return true
				}
			}
		}
	}
	return false
}
