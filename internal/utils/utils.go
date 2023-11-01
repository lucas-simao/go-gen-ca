package utils

import (
	"regexp"
	"strings"
	"text/template"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToTitle(s string) string {
	if len(s) > 0 {
		return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
	}
	return s
}

func ToLower(s string) string {
	return strings.ToLower(s)
}

func ToUpper(s string) string {
	return strings.ToUpper(s)
}

func InitTemplate(name, file string) (*template.Template, error) {
	return template.New(name).Funcs(template.FuncMap{
		"ToTitle": ToTitle,
		"ToLower": ToLower,
	}).Parse(file)
}

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
