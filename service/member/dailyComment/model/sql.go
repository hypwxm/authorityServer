package model

import (
	"babygrow/DB/pgsql"
	"fmt"
	"io/ioutil"
	"strings"
)

const table_name = "g_member_baby_grow"

func GetSqlFile() ([]byte, error) {
	b, err := ioutil.ReadFile("scheme.sql")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func insertSql() string {
	return fmt.Sprintf("insert into %s (createtime, isdelete, disabled, id, weight, height, diary, mood, temperature, weather, health, date, user_id, baby_id) select :createtime, :isdelete, :disabled, :id, :weight, :height, :diary, :mood, :temperature, :weather, :health, :date, :user_id, :baby_id returning id", table_name)

}

func listSql(query *Query) (whereSql string, fullSql string) {
	var selectSql = fmt.Sprintf(`SELECT 
				%[1]s.*,
				COALESCE(%[2]s.role_name, '') as user_role_name,
				COALESCE(%[3]s.realname, '') as user_realname,
				COALESCE(%[3]s.account, '') as user_account,
				COALESCE(%[3]s.phone, '') as user_phone,
				COALESCE(%[3]s.nickname, '') as user_nickname
				FROM %[1]s left join %[2]s on %[2]s.baby_id=%[1]s.baby_id and %[2]s.user_id=%[1]s.user_id left join %[3]s on %[1]s.user_id=%[3]s.id WHERE 1=1 `, table_name, "g_member_baby_relation", "g_member")
	whereSql = pgsql.BaseWhere(query.BaseQuery, table_name)
	if strings.TrimSpace(query.Keywords) != "" {
		// whereSql = whereSql + fmt.Sprintf(" and (%[1]s.title like '%%%[2]s%%' or %[1]s.intro like '%%%[2]s%%' or %[1]s.content like '%%%[2]s%%')", table_name, query.Keywords)
	}

	if query.UserId != "" {
		whereSql = whereSql + fmt.Sprintf(" and %[1]s.user_id=:user_id ", table_name)
	}
	if query.DiaryId != "" {
		whereSql = whereSql + fmt.Sprintf(" and %[1]s.dairy_id=:dairy_id ", table_name)
	}
	if len(query.DiaryIds) > 0 {
		whereSql = whereSql + fmt.Sprintf(" and %[1]s.dairy_id=any(:dairy_ids) ", table_name)
	}

	if query.BabyId != "" {
		whereSql = whereSql + fmt.Sprintf(" and %[1]s.baby_id=:baby_id ", table_name)
	}

	// if query.OrderBy == "" {
	// 	query.OrderBy = "sort asc"
	// } else {
	// 	query.OrderBy = "createtime desc"
	// }
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
	updateSql = updateSql + " ,weight=:weight"
	updateSql = updateSql + " ,height=:height"
	updateSql = updateSql + " ,diary=:diary"
	updateSql = updateSql + " ,weather=:weather"
	updateSql = updateSql + " ,temperature=:temperature"
	updateSql = updateSql + " ,health=:health"
	updateSql = updateSql + " ,mood=:mood"
	updateSql = updateSql + " ,date=:date"

	return fmt.Sprintf("update %s set updatetime=:updatetime %s where id=:id and isdelete=false", table_name, updateSql)
}

func delSql() string {
	return fmt.Sprintf("update %s set isdelete=true where id=any(:ids)", table_name)
}
