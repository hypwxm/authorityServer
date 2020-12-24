package model

import (
	"fmt"
)

const role_table_name = "g_admin_user_role"

func roleInsertSql() string {
	return fmt.Sprintf(`insert into %s 
	user_id, org_id
	select :user_id, :org_id`,
		org_table_name,
	)
}

func roleDelSql() string {
	return fmt.Sprintf("delete from %s where user_id=:user_id", org_table_name)
}
