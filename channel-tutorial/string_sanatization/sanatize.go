package stringsanatization

import "strings"

func Sanatize(data string) string {
	data = strings.TrimSpace(data)
	data = strings.ReplaceAll(data, "  ", " ")
	data = strings.ReplaceAll(data, "\n", " ")
	data = strings.ReplaceAll(data, "\r", " ")
	return strings.ReplaceAll(data, "\t", " ")
}
