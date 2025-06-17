package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

/*
Unpack Функция распаковывает строку и возвращает повторяющиеся символы.
*/
func Unpack(s string) (r string, err error) {
	if _, err := strconv.Atoi(s); err == nil {
		return r, errors.New("invalid string")
	}

	var prev rune
	var escaped bool
	var b strings.Builder

	for _, char := range s {
		// Является ли символ числом.
		if unicode.IsNumber(char) && !escaped {
			m := int(char - '0')
			r := strings.Repeat(string(prev), m-1)

			// Пишем в билдер новую строку, состоящую из копий символов, которые перебираем сейчас.
			b.WriteString(r)
		} else {
			// Поддержка экранирования строк.
			escaped = string(char) == "\\" && string(prev) != "\\"

			// Если пропуск не нужен, то пишем такой символ в билдер.
			if !escaped {
				b.WriteRune(char)
			}

			// Предыдущим становится текущий символ.
			prev = char
		}
	}

	// Возвращаем построенный билдер.
	return b.String(), err
}
