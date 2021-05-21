package utils

import (
	"fmt"
	"strconv"
)

func MapToInt64(m interface{}) int64 {
	mappedInt, _ := strconv.ParseInt(fmt.Sprint(m), 10, 64)
	return mappedInt
}

func TrimLeftChar(s string) string {
	for i := range s {
		if i > 0 {
			return s[i:]
		}
	}
	return s[:0]
}
