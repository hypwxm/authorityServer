package model

import (
	"babygrow/DB/pgsql"
	"fmt"
	"strings"
)

const table_name = "g_admin_role"

func insertSql() string {
	return fmt.Sprintf(`insert into %s
		(createtime, isdelete, disabled, id, name, intro, org_id) 
		select :createtime, :isdelete, :disabled, :id, :name, :intro, :org_id returning id`,
		table_name)

}

func listSql(query *Query) (whereSql string, fullSql string) {
	var selectSql = fmt.Sprintf(`SELECT 
				%[1]s.*
				FROM %[1]s WHERE 1=1 `, table_name)
	whereSql = pgsql.BaseWhere(query.BaseQuery)
	if strings.TrimSpace(query.Keywords) != "" {
		whereSql = whereSql + fmt.Sprintf(" and (%[1]s.name like '%%"+query.Keywords+"%%' or %[1]s.intro like '%%"+query.Keywords+"%%')", table_name)
	}
	if query.OrgId != "" {
		whereSql = whereSql + fmt.Sprintf(" and (%[1]s.org_id=:org_id)", table_name)
	}

	optionSql := pgsql.BaseOption(query.BaseQuery)
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
				where id=:id and isdelete=false`,
		table_name)
}

func updateSql() string {
	var updateSql = ""
	updateSql = updateSql + " ,name=:name"
	updateSql = updateSql + " ,intro=:intro"

	return fmt.Sprintf("update %s set updatetime=:updatetime %s where id=:id and isdelete=false", table_name, updateSql)
}

func hasUser() string {
	return fmt.Sprintf("select count(role_id) from %s where role_id=any(:ids)", "g_admin_user_role")
}

func delSql() string {
	return fmt.Sprintf("update %s set isdelete=true where id=any(:ids)", table_name)
}

func toggleSql() string {
	return fmt.Sprintf("update %s set disabled=:disabled where id=:id and isdelete=false", table_name)
}
