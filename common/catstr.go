package common

import (
	"bytes"
)

/*
CombineArgsWithBlank is used to cat the elements of a string array with a blank
Params: tty bool, command string, args []string
Return: error
*/
func CombineArgsWithBlank(arr []string) string {
	buf := bytes.Buffer{}
	buf.WriteString(`"`)
	buf.WriteString(arr[0])
	for _, str := range arr[1:] {
		buf.WriteString(" ")
		buf.WriteString(str)
	}
	buf.WriteString(`"`)
	resultStr := buf.String()
	return resultStr
}
