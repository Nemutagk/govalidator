package govalidator

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

type StructOptions struct {
	Tag string
}

func ValidateStruct[T any](s T, inputs []Input, customeallErrors map[string]string, models map[string]func(data any, payload map[string]any, opts *[]string) (bool, string), opts ...StructOptions) (T, error) {
	var zero T

	tag := "json"
	if len(opts) > 0 && opts[0].Tag != "" {
		tag = opts[0].Tag
	}

	body, err := structToMap(s, tag)
	if err != nil {
		return zero, fmt.Errorf("error convirtiendo struct a map: %w", err)
	}

	log.Printf("body: %+v", body)
	safePayload, err := ValidateRequest(body, inputs, customeallErrors, models)
	if err != nil {
		return zero, err
	}

	result, err := ConvertPayload[T](safePayload)
	if err != nil {
		return zero, fmt.Errorf("error convirtiendo payload a struct: %w", err)
	}

	return result, nil
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

		key, omitempty := fieldNameAndOpts(field, tag)
		if key == "-" {
			continue
		}
		if omitempty && isEmptyValue(fieldVal) {
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

func fieldNameAndOpts(field reflect.StructField, tag string) (name string, omitempty bool) {
	tagVal := field.Tag.Get(tag)
	if tagVal == "" {
		return field.Name, false
	}
	parts := strings.Split(tagVal, ",")
	name = parts[0]
	if name == "" {
		name = field.Name
	}
	for _, opt := range parts[1:] {
		if opt == "omitempty" {
			omitempty = true
		}
	}
	return name, omitempty
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	case reflect.Slice, reflect.Map, reflect.Array:
		return v.Len() == 0
	default:
		return v.IsZero()
	}
}
