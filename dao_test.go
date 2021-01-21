package code_gen

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestAddFunction(t *testing.T) {
	f, _ := os.Open("type.txt")

	defer f.Close()

	bytes, _ := ioutil.ReadAll(f)

	tokens, _ := Tokenize(string(bytes))

	typeDescList, err := Parse(tokens)

	if err != nil {
		t.Fatalf("parse failed: %s", err.Error())
	}

	for _, t := range typeDescList {
		fmt.Println(AddFunction("table", t))
	}
}
