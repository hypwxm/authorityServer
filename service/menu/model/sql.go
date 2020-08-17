package model

import (
	"fmt"
	"babygrowing/DB/pgsql"
)

const table_name = "wb_settings_menu"

func insertSql() string {
	return fmt.Sprintf("insert into %[1]s (createtime, isdelete, disabled, id, name, path, parent_id) select :createtime, :isdelete, :disabled, :id, :name, :path, :parent_id where not exists(select 1 from %[1]s where path=:path and isdelete=false) returning id", table_name)

}

func listSql(query *Query) (fullSql string) {
	if query.RoleId == "" {
		var selectSql = fmt.Sprintf(`SELECT 
				%[1]s.id,
				%[1]s.name,
				%[1]s.path,
				%[1]s.parent_id,
				%[1]s.disabled
				FROM %[1]s where 1=1 `, table_name)
		whereSql := pgsql.BaseWhere(query.BaseQuery, table_name)
		optionSql := pgsql.BaseOption(query.BaseQuery, table_name)
		return selectSql + whereSql + optionSql
	} else {
		var selectSql = fmt.Sprintf(`SELECT 
				%[2]s.id,
				%[2]s.parent_id,
				%[2]s.name,
				%[2]s.path,
				%[2]s.disabled
				FROM %[1]s inner join %[2]s on %[1]s.menu_id=%[2]s.id WHERE 1=1 `, "wb_admin_role_menu_permission", table_name)
		whereSql := pgsql.BaseWhere(query.BaseQuery, table_name)
		whereSql = whereSql + fmt.Sprintf(" and %s.role_id=:role_id", "wb_admin_role_menu_permission")

		optionSql := pgsql.BaseOption(query.BaseQuery, table_name)
		return selectSql + whereSql + optionSql
	}
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
	return fmt.Sprintf("update %s set updatetime=:updatetime %s where id=:id and isdelete=false and not exists (select 1 from %[1]s where id<>:id and isdelete=false and path=:path)", table_name, updateSql)
}

func delSql() string {
	return fmt.Sprintf("update %s set isdelete=true where id=any(:ids)", table_name)
}

func toggleSql() string {
	return fmt.Sprintf("update %s set disabled=:disabled where id=:id and isdelete=false", table_name)
}
