package code_gen

import (
	"unicode"
)

// TODO ddl => go type => dao code

func find(runes []rune, start int, tar rune) int {
	for i := start + 1; i < len(runes); i++ {
		if tar == runes[i] {
			return i
		}
	}
	return -1
}

func TokenizeDDL(typeDec string) ([]Token, error) {

	bufStrList := NewBufferedStrList()

	runes := []rune(typeDec)

	for i := 0; i < len(runes); i++ {
		r := runes[i]

		if unicode.IsSpace(r) {
			bufStrList.FlushBuffer()
		} else if r == '`' {
			bufStrList.FlushBuffer()
			j := find(runes, i, r)
			_ = bufStrList.AppendToList(string(runes[i+1 : j]))
			i = j
		} else if r == ',' || r == '(' || r == ')' {
			bufStrList.FlushBuffer()
			_ = bufStrList.AppendToList(string(r))
		}
	}

	return bufStrList.strList, nil
}
