package code_gen

import (
	"strconv"
	"strings"
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

// just remove the prefix
func parseStructName(tableName string) string {
	i := strings.Index(tableName, "_")

	return tableName[i+1:]
}

func skipWords(tokens []Token, i int) int {

	for i < len(tokens) {
		switch tokens[i] {
		case "not":
			i += 2
		case "default":
			i += 2
		case "unsigned":
			i++
		default:
			return i
		}
	}

	return i
}

func ParseDDL(tokens []Token) *StructDesc {

	structDesc := &StructDesc{
		TName:  "",
		Fields: []*FieldDesc{},
	}

	tokens = tokens[2:] // skip create table

	structDesc.TName = Token(Snake2Camel(parseStructName(tokens[0].String())))

	for i := 2; i <= len(tokens) && tokens[i] != ")"; i++ {

		if strings.ToLower(tokens[i].String()) == "primary" {
			break
		}

		columnName := tokens[i]
		i++
		columnType := tokens[i]
		i++
		length := int64(0)

		comment := ""

		if i < len(tokens) && tokens[i] == "(" {
			i++
			length, _ = strconv.ParseInt(tokens[i].String(), 10, 64)
			i += 2
		}

		i = skipWords(tokens, i)

		if i < len(tokens) && tokens[i] == "comment" {
			i += 2 // skip single quote
			comment = tokens[i].String()
			i += 2
		}

		fieldType := "string"

		if columnType == "varchar" {
			fieldType = "string"
		} else if columnType == "tinyint" {
			fieldType = "bool"
		} else if columnType == "bigint" {
			if length == 20 {
				fieldType = "int64"
			} else if length == 8 {
				fieldType = "int32"
			} else {
				fieldType = "int"
			}
		}

		fieldDesc := &FieldDesc{
			OrigFName:      columnName,
			FName:          Token(Snake2Camel(columnName.String())),
			FType:          Token(fieldType),
			IsPointer:      false,
			IsSlice:        false,
			IsPrimitive:    true,
			IsMap:          false,
			SlashedComment: comment,
		}

		structDesc.Fields = append(structDesc.Fields, fieldDesc)
	}
	return structDesc
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
		} else if r == '\'' {
			// get all letters between single quotes
			bufStrList.FlushBuffer()
			_ = bufStrList.AppendToList(string(r))
			j := i + 1
			for ; j < len(runes) && runes[j] != '\''; j++ {
			}
			bufStrList.AppendToList(string(runes[i+1 : j]))
			bufStrList.AppendToList("'")
			i = j
		} else if r == '#' {
			bufStrList.FlushBuffer()
			for runes[i] != '\n' {
				i++
			}
		} else {
			bufStrList.AppendToBuffer(string(r))
		}
	}

	return bufStrList.strList, nil
}
