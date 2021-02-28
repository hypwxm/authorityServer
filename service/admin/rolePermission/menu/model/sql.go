package model

import (
	"babygrow/DB/pgsql"
	"fmt"
)

const table_name = "g_role_menu"

func saveSql() string {
	return fmt.Sprintf(`insert into %s
		(createtime, isdelete, disabled, id, role_id, menu_id) 
		select :createtime, :isdelete, :disabled, :id, :role_id, :menu_id`,
		table_name)

}

func listSql(query *Query) (fullSql string) {
	var selectSql = fmt.Sprintf(`SELECT 
				%[1]s.createtime,
				%[1]s.updatetime,
				%[1]s.role_id,
				%[1]s.menu_id,
				%[2]s.parent_id,
				%[2]s.name,
				%[2]s.path,
				%[2]s.icon
				FROM %[1]s inner join %[2]s on %[1]s.menu_id=%[2]s.id WHERE 1=1 `, table_name, "g_menu")
	whereSql := pgsql.BaseWhere(query.BaseQuery, table_name)
	whereSql = whereSql + fmt.Sprintf(" and %[1]s.role_id=any(:role_ids)", table_name)

	optionSql := pgsql.BaseOption(query.BaseQuery, table_name)
	return selectSql + whereSql + optionSql
}

func deleteSql() string {
	return fmt.Sprintf("delete from %s where role_id=:role_id and menu_id=any(:menu_ids)", table_name)
}
