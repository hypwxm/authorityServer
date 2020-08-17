package model

import (
	"fmt"
	"babygrowing/DB/pgsql"
)

const table_name = "wb_admin_role_menu_permission"

func saveSql() string {
	return fmt.Sprintf(`insert into %s
		(createtime, isdelete, disabled, id, role_id, menu_id) 
		select :createtime, :isdelete, :disabled, :id, :role_id, :menu_id`,
		table_name)

}

func listSql(query *Query) (fullSql string) {
	var selectSql = fmt.Sprintf(`SELECT 
				%[1]s.id,
				%[1]s.createtime,
				%[1]s.updatetime,
				%[1]s.role_id,
				%[1]s.menu_id,
				%[2]s.parent_id,
				%[2]s.name,
				%[2]s.path
				FROM %[1]s inner join %[2]s on %[1]s.menu_id=%[2]s.id WHERE 1=1 `, table_name, "wb_settings_menu")
	whereSql := pgsql.BaseWhere(query.BaseQuery, table_name)
	whereSql = whereSql + fmt.Sprintf(" and %[1]s.role_id=:role_id", table_name)

	optionSql := pgsql.BaseOption(query.BaseQuery, table_name)
	return selectSql + whereSql + optionSql
}

func deleteSql(roleId string) string {
	return fmt.Sprintf("delete from %s where role_id='%s'", table_name, roleId)
}
