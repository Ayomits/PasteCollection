package utils

import "strconv"

func IsNumber(text string) bool {
	_, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		return false
	}
	return true
}

func Numberize(text string) int64 {
	i, err := strconv.ParseInt(text, 10, 32)
	if err != nil {
		return 0
	}
	return i
}
