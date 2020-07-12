package houseModel

import (
	"fmt"
)

const table_name = "wb_user_house"

func insertSql() string {
	return fmt.Sprintf("insert into %s (id, user_id, enums_id, option_id) select :id, :user_id, :enums_id, :option_id returning id", table_name)

}

func listSql(query *Query) string {
	return fmt.Sprintf(`SELECT 
				%[1]s.id,
				%[1]s.user_id,
				%[1]s.enums_id,
				%[1]s.option_id
				FROM %[1]s WHERE 1=1 and user_id=:user_id`, table_name)
}
