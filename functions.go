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
        "orVoid": OrVoid,
		"goType": GoType,
		"goTypeWithModule": GoTypeWithModule,
		"uppercaseToCapitalized": UppercaseToCapitalized,
		"replace": Replace,
	}
}

func Quote(input string) string {
	return fmt.Sprintf(`"%s"`, strings.ReplaceAll(input, `"`, `\"`))
}

func CapitalCamelCase(input string) string {
	return strings.ToUpper(input[:1]) + input[1:]
}

func OrVoid(input interface{}) string {
    if input == nil {
        return "void"
    }
    return input.(string)
}

func GoType(input string) string {
	if input == "string" {
		return "string"
	}
	if input == "number" {
		return "int64"
	}
	if input == "boolean" {
		return "bool"
	}
	return input
}

func GoTypeWithModule(module, input string) string {
	if input == "string" {
		return "string"
	}
	if input == "number" {
		return "int64"
	}
	if input == "boolean" {
		return "bool"
	}
	return fmt.Sprintf("%s%s", module, input)
}

func UppercaseToCapitalized(input string) string {
	return strings.ToUpper(input[:1]) + strings.ToLower(input[1:])
}

func Replace(whenFound string, replaceWith string, input string) string {
	return strings.ReplaceAll(input, whenFound, replaceWith)
}
