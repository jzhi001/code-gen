package code_gen

import "errors"

type TokenIterator struct {
	tokens []Token
	i      int
	n      int
}

func NewTokenIterator(tokens []Token) *TokenIterator {
	return &TokenIterator{
		tokens: tokens,
		i:      0,
		n:      len(tokens),
	}
}

func (it *TokenIterator) HasNext() bool {
	return it.i < it.n
}

func (it *TokenIterator) Next() (Token, error) {

	if !it.HasNext() {
		return "", errors.New("no more token")
	}

	tk := it.tokens[it.i]
	it.i++

	return tk, nil
}

func (it *TokenIterator) NextNonWhiteSpace() (Token, error) {

	var tk Token
	var err error

	for tk, err = it.Next(); err == nil && tk.isWhiteSpace(); tk, err = it.Next() {

	}

	if err != nil {
		return "", errors.New("no more token")
	}

	return tk, nil
}

func (it *TokenIterator) NextNonNewLine() (Token, error) {

	var tk Token
	var err error

	for tk, err = it.Next(); err == nil && tk == NewLine; tk, err = it.Next() {

	}

	if err != nil {
		return "", errors.New("no more token")
	}

	return tk, nil
}
