package normalize

import (
	"strings"
	"unicode"
)

func LTrim(value any, options []string) any {
	s, ok := value.(string)
	if !ok {
		return value
	}

	return strings.TrimLeftFunc(s, unicode.IsSpace)
}
