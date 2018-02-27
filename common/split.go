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

/*
SubStringN is used to return substring by the given number n.
Params: s string, n int
Return: string
*/
func SubStringN(s string, n int) string {
	if len(s) <= n {
		return s
	}
	str, rs := "", []rune(s)
	for _, r := range rs {
		if len(str) > n {
			break
		}
		str += string(r)
	}
	return str
}
