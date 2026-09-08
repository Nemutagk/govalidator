package normalize

import "strings"

func Lower(value any, options []string) any {
	s, ok := value.(string)
	if !ok {
		return value
	}

	return strings.ToLower(s)
}
