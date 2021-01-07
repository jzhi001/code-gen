package idl_conv

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestTokenize(t *testing.T) {

	f, _ := os.Open("type.txt")

	defer f.Close()

	bytes, _ := ioutil.ReadAll(f)

	tokens, err := Tokenize(string(bytes))

	if err != nil {
		t.Fatalf("tokenize failed: %s", err.Error())
	}

	for _, tk := range tokens {
		fmt.Printf("%s ", tk.Quote())
	}
	fmt.Println()
}
