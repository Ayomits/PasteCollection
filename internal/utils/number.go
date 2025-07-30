package utils

import "strconv"

func IsNumber(text string) bool {
	_, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		return false
	}
	return true
}
