package normalize

import "strings"

func TrimChar(value any, options []string) any {
	s, ok := value.(string)
	if !ok || len(options) == 0 {
		return value
	}

	return strings.Trim(s, strings.Join(options, ""))
}
