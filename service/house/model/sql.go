package model

import (
	"fmt"
	"worldbar/DB/pgsql"
)

const table_name = "wb_house"

func insertSql() string {
	return fmt.Sprintf("insert into %s (createtime, isdelete, disabled, id, name, note, icon, sort, parent_id) select :createtime, :isdelete, :disabled, :id, :name, :note, :icon, :sort, :parent_id returning id", table_name)

}

func listSql(query *Query) (whereSql string, fullSql string) {
	var selectSql = fmt.Sprintf(`SELECT 
				%[1]s.disabled,
				%[1]s.id,
				%[1]s.name,
				%[1]s.icon,
				%[1]s.sort,
				%[1]s.note,
				%[1]s.parent_id
				FROM %[1]s WHERE 1=1 `, table_name)
	whereSql = pgsql.BaseWhere(query.BaseQuery, table_name)
	query.OrderBy = "sort asc"
	optionSql := pgsql.BaseOption(query.BaseQuery, table_name)
	return whereSql, selectSql + whereSql + optionSql
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
	updateSql = updateSql + " ,sort=:sort"
	updateSql = updateSql + " ,icon=:icon"
	updateSql = updateSql + " ,note=:note"
	updateSql = updateSql + " ,parent_id=:parent_id"

	return fmt.Sprintf("update %s set updatetime=:updatetime %s where id=:id and isdelete=false", table_name, updateSql)
}

func delSql() string {
	return fmt.Sprintf("update %s set isdelete=true where id=any(:ids)", table_name)
}

func toggleSql() string {
	return fmt.Sprintf("update %s set disabled=:disabled where id=:id and isdelete=false", table_name)
}

