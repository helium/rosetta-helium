package utils

func MapToInt64(m interface{}) int64 {
	mappedInt := int64(m.(float64))
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

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
