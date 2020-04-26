package model

import (
	"fmt"
	"strings"
	"worldbar/DB/pgsql"
)

const table_name = "wb_settings_menu"

func insertSql() string {
	return fmt.Sprintf("insert into %[1]s (createtime, isdelete, disabled, id, name, path) select :createtime, :isdelete, :disabled, :id, :name, :path where not exists(select 1 from %[1]s where path=:path and isdelete='false') returning id", table_name)

}

func listSql(query *Query) (whereSql string, fullSql string) {
	var selectSql = fmt.Sprintf(`SELECT 
				%[1]s.id,
				%[1]s.name,
				%[1]s.path,
				%[1]s.disabled
				FROM %[1]s where 1=1 `, table_name)
	whereSql = pgsql.BaseWhere(query.BaseQuery, table_name)

	if query.OrderBy == "" {
		// query.OrderBy = "sort asc"
	}
	optionSql := pgsql.BaseOption(query.BaseQuery, table_name)
	return whereSql, selectSql + whereSql + optionSql
}

func countSql(whereSql ...string) string {
	return fmt.Sprintf("select count(*) from %s where 1=1 %s", table_name, strings.Join(whereSql, " "))
}

func getByIdSql() string {
	return fmt.Sprintf(`
			select 
				%[1]s.*
				from %[1]s
				where %[1]s.id=:id and %[1]s.isdelete=false`,
		table_name)
}

func updateSql() string {
	var updateSql = ""
	updateSql = updateSql + " ,name=:name"
	updateSql = updateSql + " ,path=:path"
	return fmt.Sprintf("update %s set updatetime=:updatetime %s where id=:id and isdelete=false", table_name, updateSql)
}

func delSql() string {
	return fmt.Sprintf("update %s set isdelete=true where id=any(:ids)", table_name)
}

func toggleSql() string {
	return fmt.Sprintf("update %s set disabled=:disabled where id=:id and isdelete=false", table_name)
}
