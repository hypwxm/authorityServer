package model

import (
	"babygrow/DB/pgsql"
	"fmt"
	"io/ioutil"
	"strings"
)

const table_name = "g_member_family"

func GetSqlFile() ([]byte, error) {
	b, err := ioutil.ReadFile("scheme.sql")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func insertSql() string {
	return fmt.Sprintf("insert into %s (createtime, isdelete, disabled, id, name, creator) select :createtime, :isdelete, :disabled, :id, :name, :creator returning id", table_name)

}

func listSql(query *Query) (whereSql string, fullSql string) {
	var selectSql = fmt.Sprintf(`SELECT 
				COALESCE(%[2]s.member_id, '') as member_id ,
				COALESCE(%[2]s.creator, '') as creator,
				COALESCE(%[2]s.can_invite, false) as can_invite,
				COALESCE(%[2]s.can_remove, false) as can_remove,
				COALESCE(%[2]s.can_edit, false) as can_edit,
				COALESCE(%[2]s.nickname, '') as nickname,
				COALESCE(%[2]s.role_name, '') as role_name,
				COALESCE(%[2]s.role_type, '') as role_type,
				COALESCE(%[1]s.createtime, 0) as createtime,
				COALESCE(%[1]s.name, '') as family_name,
				COALESCE(%[1]s.creator, '') as family_creator,
				COALESCE(%[1]s.label, '') as family_label,
				COALESCE(%[1]s.intro, '') as family_intro,
				COALESCE(%[1]s.id, '') as id
				FROM %[1]s left join %[2]s on %[1]s.id=%[2]s.family_id WHERE 1=1 `, table_name, "g_member_family_member")
	whereSql = pgsql.BaseWhere(query.BaseQuery, table_name)
	if strings.TrimSpace(query.Keywords) != "" {
		// whereSql = whereSql + fmt.Sprintf(" and (%[1]s.title like '%%%[2]s%%' or %[1]s.intro like '%%%[2]s%%' or %[1]s.content like '%%%[2]s%%')", table_name, query.Keywords)
	}
	whereSql = whereSql + fmt.Sprintf(" and (%[1]s.creator=:user_id or %[2]s.member_id=:user_id)", table_name, "g_member_family_member")

	if query.OrderBy == "" {
		// query.OrderBy = "sort asc"
	}
	optionSql := pgsql.BaseOption(query.BaseQuery, table_name)
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
				where %[1]s.id=:id and %[1]s.isdelete=false`,
		table_name)
}

func updateSql() string {
	var updateSql = ""
	updateSql = updateSql + " ,name=:name"

	return fmt.Sprintf("update %s set updatetime=:updatetime %s where id=:id and isdelete=false", table_name, updateSql)
}

func delSql() string {
	return fmt.Sprintf("update %s set isdelete=true where id=any(:ids)", table_name)
}

func toggleSql() string {
	return fmt.Sprintf("update %s set disabled=:disabled where id=:id and isdelete=false", table_name)
}
