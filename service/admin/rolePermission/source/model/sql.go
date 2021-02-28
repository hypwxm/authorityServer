package model

import (
	"babygrow/DB/pgsql"
	"fmt"
)

const table_name = "wb_admin_role_source_permission"

func saveSql() string {
	return fmt.Sprintf(`insert into %s
		(createtime, isdelete, disabled, id, role_id, source_id) 
		select :createtime, :isdelete, :disabled, :id, :role_id, :source_id`,
		table_name)

}

func listSql(query *Query) (fullSql string) {
	var selectSql = fmt.Sprintf(`SELECT 
				%[1]s.id,
				%[1]s.createtime,
				%[1]s.updatetime,
				%[1]s.role_id,
				%[1]s.source_id
				FROM %[1]s WHERE 1=1 `, table_name)
	whereSql := pgsql.BaseWhere(query.BaseQuery)
	whereSql = whereSql + fmt.Sprintf(" and role_id=:role_id")

	optionSql := pgsql.BaseOption(query.BaseQuery)
	return selectSql + whereSql + optionSql
}

func deleteSql(roleId string) string {
	return fmt.Sprintf("delete from %s where role_id='%s'", table_name, roleId)
}
