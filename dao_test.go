package code_gen

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

var tokens []Token

func init() {
	f, _ := os.Open("type.txt")

	defer f.Close()

	bytes, _ := ioutil.ReadAll(f)

	tokens, _ = Tokenize(string(bytes))
}

func TestAddFunction(t *testing.T) {

	typeDescList, err := Parse(tokens)

	if err != nil {
		t.Fatalf("parse failed: %s", err.Error())
	}

	for _, t := range typeDescList {
		fmt.Println(AddFunctionImpl("table", t))
	}
}

func TestListFunction(t *testing.T) {

	typeDescList, err := Parse(tokens)

	if err != nil {
		t.Fatalf("parse failed: %s", err.Error())
	}

	for _, t := range typeDescList {
		fmt.Println(ListFunctionImpl("table", t, []string{"GiftId", "CreatedAt"}))
	}
}

func TestUpdateByIdFunc(t *testing.T) {

	typeDescList, err := Parse(tokens)

	if err != nil {
		t.Fatalf("parse failed: %s", err.Error())
	}

	for _, t := range typeDescList {
		fmt.Println(UpdateByIdFuncImpl("table", t))
	}
}

func TestDaoCode(t *testing.T) {

	typeDescList, err := Parse(tokens)

	if err != nil {
		t.Fatalf("parse failed: %s", err.Error())
	}

	for _, t := range typeDescList {
		fmt.Println(DaoCode("table", "jzhi001_book", t, [][]string{{"Uid"}, {"Uid", "DeliverSuccess"}}))
	}
}

func TestWithDDL(t *testing.T) {

	f, _ := os.Open("example.sql")

	defer f.Close()

	bytes, _ := ioutil.ReadAll(f)

	ddl := string(bytes)

	tokens, _ := TokenizeDDL(ddl)

	structDesc, tableName := ParseDDL(tokens)

	fmt.Println(DaoCode("table", tableName, structDesc, [][]string{}))

}
