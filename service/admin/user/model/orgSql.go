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

func orgDelSql() string {
	return fmt.Sprintf("delete from %s where user_id=:user_id", org_table_name)
}
