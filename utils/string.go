package utils

import (
	"regexp"
	"strings"
)

func ConvertToUppercase(str string) string {
	str = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(str, "_")
	return strings.ToUpper(str)
}
