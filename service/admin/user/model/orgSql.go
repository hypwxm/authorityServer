package model

import (
	"fmt"
)

const org_table_name = "g_admin_user_org"

func orgInsertSql() string {
	return fmt.Sprintf(`insert into %s 
	user_id, org_id
	select :user_id, :org_id`,
		org_table_name,
	)
}

func orgListSql() (whereSql string) {
	var selectSql = fmt.Sprintf(`SELECT 
				%[1]s.*,
				%[2]s.*
				FROM %[1]s left join %[2]s on %[1]s.org_id=%[2]s.id WHERE 1=1 and %[2]s.isdelete=false and %[1]s.user_id=any(:user_ids) `, table_name, "g_admin_org")

	return selectSql
}

func orgDelSql() string {
	return fmt.Sprintf("delete from %s where user_id=:user_id", org_table_name)
}
