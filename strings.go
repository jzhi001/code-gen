package code_gen

import (
	"strings"
	"unicode"
)

func Camel2Snake(s string) string {

	cit := NewCharIterator(s)
	buf := NewBufferedStrList()

	for cit.HasNext() {

		r, _ := cit.NextRune()

		if unicode.IsUpper(r) {
			buf.FlushBuffer()
			r = unicode.ToLower(r)
		}
		buf.AppendToBuffer(string(r))
	}
	buf.FlushBuffer()

	var list []string

	for _, tk := range buf.GetList() {
		list = append(list, tk.String())
	}

	return strings.Join(list, "_")
}

func Snake2Camel(s string) string {

	cit := NewCharIterator(s)
	buf := NewBufferedStrList()

	for cit.HasNext() {

		r, _ := cit.NextRune()

		if r == '_' {
			buf.FlushBuffer()
		} else {
			buf.AppendToBuffer(string(r))
		}
	}
	buf.FlushBuffer()

	var list []string

	for _, tk := range buf.GetList() {
		list = append(list, strings.Title(tk.String()))
	}

	return strings.Join(list, "")
}

func UnTitle(s string) string {
	runes := []rune(s)

	if len(runes) == 0 {
		return s
	}

	if unicode.IsUpper(runes[0]) {
		runes[0] = unicode.ToLower(runes[0])
	}

	return string(runes)
}
