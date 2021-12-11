package codegenerator

import (
	"fmt"
	"strings"
	"text/template"
)

func AuxilirayFunctions() template.FuncMap {
	return template.FuncMap{
		"quote": Quote,
		"capitalCamelCase": CapitalCamelCase,
	}
}

func Quote(input string) string {
	return fmt.Sprintf(`"%s"`, strings.ReplaceAll(input, `"`, `\"`))
}

func CapitalCamelCase(input string) string {
	return strings.ToUpper(input[:1]) + input[1:]
}
