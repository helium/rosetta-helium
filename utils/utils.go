package utils

import (
	"fmt"
	"strconv"
)

func MapToInt64(m interface{}) int64 {
	mappedInt, _ := strconv.ParseInt(fmt.Sprint(m), 10, 64)
	return mappedInt
}
