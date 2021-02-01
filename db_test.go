package code_gen

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestTokenizeDDL(t *testing.T) {
	f, _ := os.Open("example.sql")

	defer f.Close()

	bytes, _ := ioutil.ReadAll(f)

	ddl := string(bytes)

	tokens, _ := TokenizeDDL(ddl)

	for _, token := range tokens {
		fmt.Println(token)
	}
}

func TestParseDDL(t *testing.T) {
	f, _ := os.Open("example.sql")

	defer f.Close()

	bytes, _ := ioutil.ReadAll(f)

	ddl := string(bytes)

	tokens, _ := TokenizeDDL(ddl)

	structDesc, _ := ParseDDL(tokens)

	println(structDesc.TName)

	for _, field := range structDesc.Fields {
		fmt.Printf("%s %s %s\n", field.FName, field.FType, field.SlashedComment)
	}
}
