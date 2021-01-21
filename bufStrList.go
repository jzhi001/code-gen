package code_gen

import "errors"

type BufferedTokenList struct {
	strList []Token
	cur     string
}

func NewBufferedStrList() *BufferedTokenList {
	return &BufferedTokenList{
		strList: []Token{},
		cur:     "",
	}
}

func (lst *BufferedTokenList) AppendToBuffer(s string) {

	lst.cur += s
}

func (lst *BufferedTokenList) FlushBuffer() {

	if len(lst.cur) == 0 {
		return
	}

	lst.strList = append(lst.strList, Token(lst.cur))
	lst.cur = ""
}

func (lst *BufferedTokenList) GetList() []Token {
	return lst.strList
}

func (lst *BufferedTokenList) AppendToList(tar string) error {

	if len(lst.cur) > 0 {
		return errors.New("cannot append to list while buffer is not empty")
	}

	lst.strList = append(lst.strList, Token(tar))

	return nil
}
