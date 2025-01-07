package stringsanatization

import (
	"strings"
	"time"
)

func Sanatize(data string) string {
	data = strings.TrimSpace(data)
	data = strings.ReplaceAll(data, "  ", " ")
	data = strings.ReplaceAll(data, "\n", " ")
	data = strings.ReplaceAll(data, "\r", " ")
	time.Sleep(500 * time.Millisecond)
	return strings.ReplaceAll(data, "\t", " ")
}
