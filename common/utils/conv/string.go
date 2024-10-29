package conv

import (
	"unicode"
	"unicode/utf8"
)

func ToCamelCase(snake_case string) string {
	if snake_case == "" {
		return ""
	}

	var bytes []byte
	toUpper := false
	for i := 0; i < len(snake_case); i++ {
		c := snake_case[i]
		switch c < utf8.RuneSelf {
		case toUpper:
			bytes = append(bytes, byte(unicode.ToUpper(rune(c))))
			toUpper = false
		case unicode.IsSpace(rune(c)):
			toUpper = true
		default:
			bytes = append(bytes, byte(unicode.ToLower(rune(c))))
		}
	}
	return B2S(bytes)
}

func TrimSpacePrefix(s string) string {
	if s == "" {
		return ""
	}

	var bytes []byte
	for i := 0; i < len(s); i++ {
		if !unicode.IsSpace(rune(s[i])) {
			bytes = append(bytes, s[i:]...)
			break
		}
	}
	if bytes == nil {
		return ""
	}
	return B2S(bytes)
}
