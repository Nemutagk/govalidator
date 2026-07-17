package helper

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/google/uuid"
)

func GenerateUuid() string {
	buil, err := uuid.NewV7()
	if err != nil {
		return ""
	}

	return buil.String()
}

func PrettyPrint(data any) {
	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Error formatting JSON:", err)
		return
	}

	fmt.Println(string(prettyJSON))
}

func SliceContains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func IsEmpty(v any) bool {
	if v == nil {
		return true
	}

	rv := reflect.ValueOf(v)

	// Desenvuelve punteros e interfaces
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return true
		}
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.String:
		return strings.TrimSpace(rv.String()) == ""
	case reflect.Slice, reflect.Array, reflect.Map:
		return rv.Len() == 0
	case reflect.Bool:
		return false // normalmente false NO se considera "vacío"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return false // si quieres, aquí podrías considerar 0 como vacío
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return false
	case reflect.Float32, reflect.Float64:
		return false
	case reflect.Struct:
		// opcional: comparar con zero-value del struct
		zero := reflect.Zero(rv.Type())
		return reflect.DeepEqual(rv.Interface(), zero.Interface())
	default:
		return false
	}
}
