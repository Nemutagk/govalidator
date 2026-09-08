package normalize

import (
	"strings"
	"unicode"
)

func RTrim(value any, options []string) any {
	s, ok := value.(string)
	if !ok {
		return value
	}

	return strings.TrimRightFunc(s, unicode.IsSpace)
}
