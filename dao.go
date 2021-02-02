package code_gen

import (
	"strings"
)

// generate DAO code

const tab = "\t"

func clusterName(table string) string {
	ans := ""

	for i := 0; i < len(table) && table[i] != '_'; i++ {
		ans += string(table[i])
	}

	return strings.ToUpper(ans)

}

func DaoCode(tableName string, typeDesc *StructDesc, queryCriteria [][]string) string {

	ans := "package dao\n\n"
	ans += "import (\n\"context\"\n\"fmt\"\n)\n"

	tableVar := Snake2Camel(tableName) + "Table"

	ans += "const cluster = \"" + clusterName(tableName) + "\"\n"
	ans += "const " + tableVar + " = \"" + tableName + "\"\n"
	ans += "const partition = 500"

	ans += typeDesc.String() + "\n"

	ans += "type " + typeDesc.TName.String() + "Dao interface{\n"

	ans += AddFunction(tableVar, typeDesc)
	ans += BatchAddFunction(tableName, typeDesc)
	ans += UpdateByIdFunc(tableVar, typeDesc)
	ans += ListFunction(tableVar, typeDesc, []string{})

	for _, criterion := range queryCriteria {
		ans += ListFunction(tableVar, typeDesc, criterion)
	}

	ans += "}\n\n"

	ans += "type " + typeDesc.TName.String() + "DaoImpl "
	ans += `struct{
	Db    *util.DB
}
`
	ans += "var Default" + typeDesc.TName.String() + "Dao " + typeDesc.TName.String() + "Dao\n"
	ans += "var _ " + typeDesc.TName.String() + "Dao = (*" + typeDesc.TName.String() + "DaoImpl)(nil)\n"

	ans += "func init() {\n"
	ans += "db := util.NewDB(cluster, " + tableVar + ")\n"
	ans += "Default" + typeDesc.TName.String() + "Dao = &" + typeDesc.TName.String() + "DaoImpl{\n"
	ans += "Db: db,\n"
	ans += "}}\n"

	ans += AddFunctionImpl(tableVar, typeDesc)
	ans += BatchAddFunctionImpl(tableVar, typeDesc)
	ans += UpdateByIdFuncImpl(tableVar, typeDesc)
	ans += ListFunctionImpl(tableVar, typeDesc, []string{})

	for _, criterion := range queryCriteria {
		ans += ListFunctionImpl(tableVar, typeDesc, criterion)
	}

	return ans
}

func daoImplReceiver(typeDesc *StructDesc) string {
	return "(dao *" + typeDesc.TName.String() + "DaoImpl)"
}

func addFuncNameAndParam(typeDesc *StructDesc) string {
	return "Add(ctx context.Context, model *" + typeDesc.TName.String() + ") error"
}

func columnList(typeDesc *StructDesc) string {

	var columns []string

	for _, f := range typeDesc.Fields {
		columns = append(columns, Camel2Snake(f.FName.String()))
	}

	return strings.Join(columns, ", ")
}

func placeHolderList(typeDesc *StructDesc) string {
	questionMarks := strings.Repeat("?", len(typeDesc.Fields))
	questionMarkList := strings.Split(questionMarks, "")
	return strings.Join(questionMarkList, ", ")
}

func listFuncName(typeDesc *StructDesc, criteriaList []string) string {
	fun := "List"

	if len(criteriaList) == 0 {
		return fun
	}

	fun += "By"

	for i, criteria := range criteriaList {
		field, err := typeDesc.GetField(criteria)
		if err != nil {
			panic("invalid criteria")
		}
		if i > 0 && len(criteriaList) < 3 {
			fun += "And"
		}
		fun += field.FName.String()
	}

	return fun
}

func listOfPt(t string) string {
	return "[]*" + t
}

func mapToUntitled(strList []string) []string {
	var ans []string

	for _, s := range strList {
		ans = append(ans, UnTitle(s))
	}

	return ans
}

func ListFunction(tableVar string, typeDesc *StructDesc, criteriaList []string) string {
	fun := listFuncName(typeDesc, criteriaList) + "(ctx context.Context"

	for _, criteria := range criteriaList {

		field, err := typeDesc.GetField(criteria)
		if err != nil {
			panic("invalid criteria")
		}

		fun += ", " + UnTitle(field.FName.String()) + " " + field.FType.String()
	}

	fun += ") (" + listOfPt(typeDesc.TName.String()) + ", error)\n"

	return fun
}

func ListFunctionImpl(tableVar string, typeDesc *StructDesc, criteriaList []string) string {
	fun := "func " + daoImplReceiver(typeDesc) + " " + listFuncName(typeDesc, criteriaList) + "(ctx context.Context"

	for _, criteria := range criteriaList {

		field, err := typeDesc.GetField(criteria)
		if err != nil {
			panic("invalid criteria")
		}

		fun += ", " + UnTitle(field.FName.String()) + " " + field.FType.String()
	}

	fun += ") (" + listOfPt(typeDesc.TName.String()) + ", error) {\n"

	fun += tab + `sql := fmt.Sprintf("SELECT ` + columnList(typeDesc) + ` FROM %s`

	if len(criteriaList) > 0 {
		fun += " WHERE"
		for i, criteria := range criteriaList {

			field, err := typeDesc.GetField(criteria)
			if err != nil {
				panic("invalid criteria")
			}
			if i > 0 {
				fun += " AND"
			}
			fun += " " + Camel2Snake(field.FName.String()) + " = ?"
		}
	}
	fun += `", ` + tableVar + ")\n"

	fun += "\tr, err := dao.Db.QueryContext(ctx, sql, " + strings.Join(mapToUntitled(criteriaList), ", ") + ")\n"

	// TODO auto tab
	fun += "\tif err != nil {\n \t\t return nil, err\n\t}\n"
	fun += tab + "defer r.Close()\n\n"
	fun += tab + "var list " + listOfPt(typeDesc.TName.String()) + "\n"

	fun += "for r.Next(){\n"
	fun += UnTitle(typeDesc.TName.String()) + " := " + typeDesc.TName.String() + "{}\n"
	fun += "err = r.Scan("

	for i, field := range typeDesc.Fields {
		if i > 0 {
			fun += ", "
		}
		fun += "&" + UnTitle(typeDesc.TName.String()) + "." + field.FName.String()
	}
	fun += ")\n"
	fun += `if err != nil {
			return nil, err
		}
`
	fun += "list = append(list, &" + UnTitle(typeDesc.TName.String()) + ")\n"
	fun += "}\n"

	fun += "return list, nil\n"

	fun += "}\n"

	return fun
}

func AddFunction(tableVar string, typeDesc *StructDesc) string {
	return "Add(ctx context.Context, model *" + typeDesc.TName.String() + ") error\n"
}

func BatchAddFunction(tableVar string, typeDesc *StructDesc) string {
	return "BatchAdd(ctx context.Context, models []*" + typeDesc.TName.String() + ") error\n"
}

func AddFunctionImpl(tableVar string, typeDesc *StructDesc) string {

	fun := "func " + daoImplReceiver(typeDesc) + " Add(ctx context.Context, model *" + typeDesc.TName.String() + ") error{\n"

	// TODO statement layer receiver.call(func, params).String()
	// TODO sql layer
	fun += tab + `sql := fmt.Sprintf("INSERT INTO %s (` + columnList(typeDesc) + ") VALUES (" + placeHolderList(typeDesc) + `)", ` + tableVar + ")\n"

	fun += tab + `_, err := dao.Db.ExecContext(ctx, sql`

	for _, f := range typeDesc.Fields {
		fun += ", model." + f.FName.String()
	}

	fun += ")\n"
	fun += tab + "return err\n}\n"

	return fun
}

func BatchAddFunctionImpl(tableVar string, typeDesc *StructDesc) string {

	fun := "func " + daoImplReceiver(typeDesc) + " BatchAdd(ctx context.Context, models []*" + typeDesc.TName.String() + ") error{\n"

	fun += `for i := 0; i < len(models); i += partition {` + "\n"
	// TODO statement layer receiver.call(func, params).String()
	// TODO sql layer
	fun += tab + `sql := fmt.Sprintf("INSERT INTO %s (` + columnList(typeDesc) + `) VALUES ", ` + tableVar + ")\n"

	fun += `for j, c := i, 0; j < len(models) && c < partition; j, c = j+1, c+1 {` + "\n"
	fun += "model := models[j]\n"

	fun += `if j > 0 {
			sql += ","
		}
`
	fun += `sql += "("`

	for i, field := range typeDesc.Fields {
		if i > 0 {
			fun += `+ ","`
		}
		if field.FType == "string" {
			fun += `+ "'" + model.` + field.FName.String() + ` + "'"`
		} else {
			if field.FType == "int64" {
				fun += "+ strconv.FormatInt(model." + field.FName.String() + ", 10)"
			} else {
				fun += "+ strconv.FormatInt(int64(model." + field.FName.String() + "), 10)"
			}

		}
	}
	fun += `+ ")"`

	fun += "\n}\n"

	fun += tab + `_, err := dao.Db.ExecContext(ctx, sql)` + "\n"
	fun += `if err != nil {

		}
`
	fun += "}\n"
	fun += tab + "return nil\n}\n"

	return fun
}

func UpdateByIdFunc(tableVar string, typeDesc *StructDesc) string {
	return "UpdateById(ctx context.Context, model *" + typeDesc.TName.String() + ") error\n"
}

func UpdateByIdFuncImpl(tableVar string, typeDesc *StructDesc) string {

	fun := "func " + daoImplReceiver(typeDesc) + " UpdateById(ctx context.Context, model *" + typeDesc.TName.String() + ") error{\n"

	fun += tab + `sql := fmt.Sprintf("UPDATE %s SET`

	for i, field := range typeDesc.Fields {
		// skip id
		// TODO identify id field by convention name or tag
		if i == 0 {
			continue
		}
		if i > 1 {
			fun += ", "
		}
		fun += " `" + Camel2Snake(field.FName.String()) + "` = ?"
	}
	fun += " WHERE `id` = ?\", " + tableVar + ")\n"

	fun += tab + `_, err := dao.Db.ExecContext(ctx, sql`

	for _, f := range typeDesc.Fields {
		if f.FName.String() == "Id" {
			continue
		}
		fun += ", model." + f.FName.String()
	}
	fun += ", model.Id"

	fun += ")\n"
	fun += tab + "return err\n}\n"

	return fun
}
