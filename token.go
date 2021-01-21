package code_gen

import (
	"strconv"
	"strings"
)

type Token string

const (
	NewLine       = "\n"
	BackTick      = "`"
	Space         = " "
	Tab           = "\t"
	LfCurlBrace   = "{"
	RtCurlBrace   = "}"
	LfSquareBrace = "["
	RtSquareBrace = "]"
	Slash         = "/"
	Carriage      = "\r"
	Type          = "type"
)

func (tk Token) isWhiteSpace() bool {
	switch tk {
	case Space, Tab, Carriage:
		return true
	default:
		return false
	}
}

func (tk Token) Quote() string {
	return strconv.Quote(string(tk))
}

func (tk Token) StartsWith(tar string) bool {

	return strings.HasPrefix(string(tk), tar)
}

func (tk Token) String() string {
	return string(tk)
}

func Tokenize(typeDec string) ([]Token, error) {

	bufStrList := NewBufferedStrList()
	cit := NewCharIterator(typeDec)

	for c, err := cit.NextChar(); cit.HasNext(); c, err = cit.NextChar() {

		if err != nil {
			return nil, err
		}

		switch c {
		case Space, Tab, Carriage:
			bufStrList.FlushBuffer()
		case Slash:
			bufStrList.FlushBuffer()
			err := cit.SkipUntil([]rune(NewLine)[0])
			if err != nil {
				return nil, err
			}
		case BackTick:
			bufStrList.FlushBuffer()
			err := cit.JumpTo([]rune(BackTick)[0])
			if err != nil {
				return nil, err
			}
		case LfCurlBrace, RtCurlBrace, NewLine:
			bufStrList.FlushBuffer()
			_ = bufStrList.AppendToList(c)
		default:
			bufStrList.AppendToBuffer(c)
		}
	}

	return bufStrList.strList, nil
}
