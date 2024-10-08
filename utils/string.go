package utils

import (
	"encoding/json"
	"regexp"
	"strings"
)

func ConvertToUpperCase(str string) string {
	str = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(str, "_")
	return strings.ToUpper(str)
}

func ConvertToCamelCase(str string) string {
	str = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(str, " ")
	parts := strings.Split(str, " ")
	for i := 0; i < len(parts); i++ {
		parts[i] = strings.Title(parts[i])
	}

	// Concatenate the words
	result := strings.Join(parts, "")
	return result
}

func ConvertStructToString(data interface{}) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(jsonData)
}
