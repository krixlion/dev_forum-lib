package str

import (
	"regexp"
	"strings"
)

func ToLowerSnakeCase(str string) string {
	matchAllCap := regexp.MustCompile("([a-z0-9])([A-Z])")
	matchFirstCap := regexp.MustCompile("(.)([A-Z][a-z]+)")

	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")

	return strings.ToLower(snake)
}
