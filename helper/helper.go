package helper

import (
	"encoding/json"
	"fmt"

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
