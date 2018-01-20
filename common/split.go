package common

import (
	"strings"
	"unicode"
)

/*
SplitArgs is used to split args by WhiteSpace
Params: args []string
Return: []string
*/
func SplitArgs(args []string) []string {
	argv := []string{}
	for _, arg := range args {
		argv = append(argv, strings.FieldsFunc(arg, unicode.IsSpace)...)
	}
	return argv
}
