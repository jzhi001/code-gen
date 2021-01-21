package code_gen

import "strings"

// generate DAO code

const tab = "\t"

//func DaoCode(tableVar string, typeDesc *TypeDesc, queryCriteria [][]string) string{
//
//
//}

// type CommodityDao interface {
func DaoTypeCode(typeDesc *TypeDesc) string {
	code := "type " + typeDesc.TName.String() + "Dao interface{\n"

	// add function
	code += "\tfunc " + addFuncNameAndParam(typeDesc) + "\n"

	code += "}\n"

	return code
}

func daoImplReceiver(typeDesc *TypeDesc) string {
	return "(dao *" + typeDesc.TName.String() + "DaoImpl)"
}

func addFuncNameAndParam(typeDesc *TypeDesc) string {
	return "Add(ctx context.Context, model *" + typeDesc.TName.String() + ") error"
}

func columnList(typeDesc *TypeDesc) string {

	var columns []string

	for _, f := range typeDesc.Fields {
		columns = append(columns, f.FName.String())
	}

	return strings.Join(columns, " ,")
}

func placeHolderList(typeDesc *TypeDesc) string {
	questionMarks := strings.Repeat("?", len(typeDesc.Fields))
	questionMarkList := strings.Split(questionMarks, "")
	return strings.Join(questionMarkList, ", ")
}

func AddFunction(tableVar string, typeDesc *TypeDesc) string {
	//daoType := "A"
	fun := "func " + daoImplReceiver(typeDesc) + " Add(ctx context.Context, model *" + typeDesc.TName.String() + ") error{\n"

	// TODO statement layer receiver.call(func, params).String()
	// TODO sql layer
	fun += tab + `fmt.Sprintf("INSERT INTO %s (` + columnList(typeDesc) + ") VALUES (" + placeHolderList(typeDesc) + `)", ` + tableVar + ")\n"

	fun += tab + `_, err := dao.Db.ExecContext(ctx`

	for _, f := range typeDesc.Fields {
		fun += ", model." + f.OrigFName.String()
	}

	fun += ")\n"
	fun += tab + "return err\n}\n"

	return fun
}
