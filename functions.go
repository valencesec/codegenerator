package codegenerator

import (
	"fmt"
	"strings"
	"text/template"
	"unicode"
)

func AuxilirayFunctions() template.FuncMap {
	return template.FuncMap{
		"quote":                        Quote,
		"capitalCamelCase":             CapitalCamelCase,
		"orVoid":                       OrVoid,
		"goType":                       GoType,
		"goTypeWithModule":             GoTypeWithModule,
		"uppercaseToCapitalized":       UppercaseToCapitalized,
		"replace":                      Replace,
		"split":                        Split,
		"splitIndexOf":                 SplitIndexOf,
		"splitIndexOfNegative":         SplitIndexOfNegative,
		"lowerNoUnderscore":            LowerNoUnderscore,
		"camelCaseNoUnderscore":        CamelCaseNoUnderscore,
		"capitalCamelCaseNoUnderscore": CapitalCamelCaseNoUnderscore,
		"upperSpaceToUnderscore":       UpperSpaceToUnderscore,
		"rustType":                     RustType,
		"camelCaseToLowerSnakeCase":    CamelCaseToLowerSnakeCase,
		"has":                          Has,
		"trimPrefix":                   TrimPrefix,
		"trimSuffix":                   TrimSuffix,
		"isBasicType":                  IsBasicType,
		"cutTakeBefore":                CutTakeBefore,
		"cutTakeAfter":                 CutTakeAfter,
	}
}

func Has(input interface{}, field string) bool {
	asMap, ok := input.(map[interface{}]interface{})
	if !ok {
		return false
	}
	_, has := asMap[field]
	return has
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

func IsBasicType(input string) bool {
	return input == "string" || input == "number" || input == "boolean" || input == "unknown"
}

func goType(input string) string {
	if input == "string" {
		return "string"
	}
	if input == "number" {
		return "int64"
	}
	if input == "boolean" {
		return "bool"
	}
	if input == "unknown" {
		return "interface{}"
	}
	return ""
}

func GoType(input string) string {
	converted := goType(input)
	if converted != "" {
		return converted
	}
	return input
}

func GoTypeWithModule(module, input string) string {
	converted := goType(input)
	if converted != "" {
		return converted
	}
	return fmt.Sprintf("%s%s", module, input)
}

func UppercaseToCapitalized(input string) string {
	return strings.ToUpper(input[:1]) + strings.ToLower(input[1:])
}

func Replace(whenFound string, replaceWith string, input string) string {
	return strings.ReplaceAll(input, whenFound, replaceWith)
}

func Split(seperator string, input string) []string {
	return strings.Split(input, seperator)
}

func SplitIndexOf(seperator string, lookFor string, input string) int {
	parts := strings.Split(input, seperator)
	for i, part := range parts {
		if part == lookFor {
			return i
		}
	}
	return -1
}

func SplitIndexOfNegative(seperator string, lookFor string, input string) int {
	parts := strings.Split(input, seperator)
	for i, part := range parts {
		if part == lookFor {
			return i - len(parts)
		}
	}
	return -1000000
}

func LowerNoUnderscore(input string) string {
	return strings.ToLower(strings.ReplaceAll(input, "_", ""))
}

func CapitalCamelCaseNoUnderscore(input string) string {
	parts := strings.Split(input, "_")
	for i, part := range parts {
		parts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
	}
	return strings.Join(parts, "")
}

func CamelCaseNoUnderscore(input string) string {
	capital := CapitalCamelCaseNoUnderscore(input)
	return strings.ToLower(capital[:1]) + capital[1:]
}

func UpperSpaceToUnderscore(input string) string {
	return strings.ToUpper(strings.ReplaceAll(input, " ", "_"))
}

func RustType(input string) string {
	if input == "string" {
		return "String"
	}
	if input == "number" {
		return "i64"
	}
	if input == "boolean" {
		return "bool"
	}
	return input
}

func CamelCaseToLowerSnakeCase(input string) string {
	if len(input) == 0 {
		return ""
	}
	result := strings.ToLower(input[:1])
	for i, runeValue := range input {
		if i == 0 {
			continue
		}
		if unicode.IsUpper(runeValue) {
			result += "_" + string(unicode.ToLower(runeValue))
		} else {
			result += string(runeValue)
		}
	}
	return result
}

func TrimPrefix(prefix string, input string) string {
	return strings.TrimPrefix(input, prefix)
}

func TrimSuffix(suffix string, input string) string {
	return strings.TrimSuffix(input, suffix)
}

func CutTakeBefore(seperator string, input string) string {
	before, _, ok := strings.Cut(input, seperator)
	if ok {
		return before
	} else {
		return ""
	}
}

func CutTakeAfter(seperator string, input string) string {
	_, after, ok := strings.Cut(input, seperator)
	if ok {
		return after
	} else {
		return ""
	}
}
