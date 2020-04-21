package utils

import (
	"unicode/utf8"
)

// PasswordValidtion パスワードのバリデーション判定を行う
func PasswordValidtion(password string) (result bool) {
	// 現段階では文字数制限のみ
	length := utf8.RuneCountInString(password)
	if length < 8 {
		result = false
	} else {
		result = true
	}
	return result
}
