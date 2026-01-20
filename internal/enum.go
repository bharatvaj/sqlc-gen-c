package c

import (
	"strings"
	"unicode"
)

type Constant struct {
	Type  string
	Name  string
	Value string
}

type Enum struct {
	Name    string
	Values  []string
	Comment string
}

func enumReplacer(r rune) rune {
	if strings.ContainsRune("-/:_", r) {
		return '_'
	} else if (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9') {
		return r
	} else {
		return -1
	}
}

// EnumReplace removes all non ident symbols (all but letters, numbers and
// underscore) and returns valid ident name for provided name.
func EnumReplace(value string) string {
	return strings.Map(enumReplacer, value)
}

// EnumValueName removes all non ident symbols (all but letters, numbers and
// underscore) and converts snake case ident to camel case.
func EnumValueName(value string) string {
	parts := strings.Split(EnumReplace(value), "_")
	for i, part := range parts {
		parts[i] = titleFirst(part)
	}

	return strings.Join(parts, "")
}

func titleFirst(s string) string {
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])

	return string(r)
}
