package main

import (
	"unicode"
)

// 判断是否整数
func isNumber(content string) bool {
	for _, r := range content {
		if !unicode.IsNumber(r) {
			return false
		}
	}
	return true
}
