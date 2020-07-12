package houseModel

import (
	"fmt"
)

const table_name = "wb_user_house"

func insertSql() string {
	return fmt.Sprintf("insert into %s (id, user_id, enums_id, option_id) select :id, :user_id, :enums_id, :option_id returning id", table_name)

}

func listSql(query *GetQuery) string {
	return fmt.Sprintf(`SELECT 
				%[1]s.id,
				%[1]s.user_id,
				%[1]s.enums_id,
				%[1]s.option_id,
				%[2]s.name as enums_name,
				%[3]s.name as option_name
				FROM %[1]s 
				inner join %[2]s on %[1]s.enums_id=%[2]s.id 
 				inner join %[3]s on %[1]s.option_id=%[3]s.id 
				WHERE 1=1 and %[1]s.user_id=:user_id`,
				table_name, "wb_house_enums", "wb_house_option")
}
