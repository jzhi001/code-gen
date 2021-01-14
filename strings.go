package idl_conv

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
