package orm

import (
	"bytes"
	"unicode"
)

type MapperFn func(string) string

func DbMapper(n string) string {
	var buf bytes.Buffer
	for i, it := range n {
		if i == 0 {
			buf.WriteRune(unicode.ToLower(it))
		} else if unicode.IsUpper(it) {
			buf.WriteString("_")
			buf.WriteRune(unicode.ToLower(it))
		} else {
			buf.WriteRune(it)
		}
	}
	return buf.String()
}

func SameMapper(n string) string {return n}
