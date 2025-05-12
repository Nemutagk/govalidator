package helper

import (
	"encoding/json"
	"fmt"

	"github.com/gofrs/uuid"
)

func GenerateUuid() string {
	return uuid.Must(uuid.NewV7()).String()
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
