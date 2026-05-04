package govalidator

import (
	"fmt"
	"reflect"
	"strings"
)

type StructOptions struct {
	Tag string
}

func ValidateStruct(s any, inputs []Input, customeallErrors map[string]string, models map[string]func(data any, payload map[string]any, opts *[]string) (bool, string), opts ...StructOptions) (map[string]any, error) {
	tag := "json"
	if len(opts) > 0 && opts[0].Tag != "" {
		tag = opts[0].Tag
	}

	body, err := structToMap(s, tag)
	if err != nil {
		return nil, fmt.Errorf("error convirtiendo struct a map: %w", err)
	}

	return ValidateRequest(body, inputs, customeallErrors, models)
}

func structToMap(s any, tag string) (map[string]any, error) {
	return reflectStructToMap(reflect.ValueOf(s), tag)
}

func reflectStructToMap(v reflect.Value, tag string) (map[string]any, error) {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil, nil
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("se esperaba un struct, se recibió %s", v.Kind())
	}

	result := make(map[string]any)
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldVal := v.Field(i)

		if !field.IsExported() {
			continue
		}

		key := fieldName(field, tag)
		if key == "-" {
			continue
		}

		val, err := reflectToAny(fieldVal, tag)
		if err != nil {
			return nil, err
		}
		result[key] = val
	}

	return result, nil
}

func reflectToAny(v reflect.Value, tag string) (any, error) {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil, nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		if v.Type().String() == "time.Time" {
			return v.Interface(), nil
		}
		return reflectStructToMap(v, tag)
	case reflect.Slice:
		if v.IsNil() {
			return nil, nil
		}
		// Preservar tipo de slices primitivos (no recursivos)
		elemKind := v.Type().Elem().Kind()
		if elemKind == reflect.Ptr {
			elemKind = v.Type().Elem().Elem().Kind()
		}
		if elemKind != reflect.Struct && elemKind != reflect.Interface && elemKind != reflect.Map {
			return v.Interface(), nil
		}

		result := make([]any, v.Len())
		for i := range v.Len() {
			val, err := reflectToAny(v.Index(i), tag)
			if err != nil {
				return nil, err
			}
			result[i] = val
		}
		return result, nil
	case reflect.Map:
		if v.IsNil() {
			// fallback para mapas nulos
		}
		// Preservar mapas primitivos si es necesario, o al menos convertir controladamente
		elemKind := v.Type().Elem().Kind()
		if elemKind == reflect.Ptr {
			elemKind = v.Type().Elem().Elem().Kind()
		}
		if elemKind != reflect.Struct && elemKind != reflect.Interface && elemKind != reflect.Map {
			return v.Interface(), nil
		}

		result := make(map[string]any)
		for _, k := range v.MapKeys() {
			val, err := reflectToAny(v.MapIndex(k), tag)
			if err != nil {
				return nil, err
			}
			result[fmt.Sprintf("%v", k.Interface())] = val
		}
		return result, nil
	default:
		return v.Interface(), nil
	}
}

func fieldName(field reflect.StructField, tag string) string {
	if tagVal := field.Tag.Get(tag); tagVal != "" {
		name, _, _ := strings.Cut(tagVal, ",")
		if name != "" {
			return name
		}
	}
	return field.Name
}
