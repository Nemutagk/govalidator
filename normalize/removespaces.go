package normalize

import "strings"

func RemoveSpaces(value any, options []string) any {
	s, ok := value.(string)
	if !ok {
		return value
	}

	return strings.Join(strings.Fields(s), "")
}
