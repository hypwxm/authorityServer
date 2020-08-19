package model

import (
	"babygrowing/DB/pgsql"
	"fmt"
	"strings"
)

const table_name = "g_daily"

func insertSql() string {
	return fmt.Sprintf("insert into %s (createtime, isdelete, disabled, id, weight, height, diary, user_id, baby_id) select :createtime, :isdelete, :disabled, :id, :weight, :height, :diary, :user_id, :baby_id returning id", table_name)

}

func listSql(query *Query) (whereSql string, fullSql string) {
	var selectSql = fmt.Sprintf(`SELECT 
				%[1]s.id,
				%[1]s.createtime,
				%[1]s.updatetime,
				%[1]s.weight,
				%[1]s.height,
				%[1]s.diary,
				%[1]s.user_id,
				%[1]s.baby_id
				FROM %[1]s WHERE 1=1 `, table_name)
	whereSql = pgsql.BaseWhere(query.BaseQuery, table_name)
	if strings.TrimSpace(query.Keywords) != "" {
		// whereSql = whereSql + fmt.Sprintf(" and (%[1]s.title like '%%%[2]s%%' or %[1]s.intro like '%%%[2]s%%' or %[1]s.content like '%%%[2]s%%')", table_name, query.Keywords)
	}

	if query.PublishTime > 0 {
		whereSql = whereSql + fmt.Sprintf(" and %s.publish_time<=:publish_time", table_name)
	}

	if query.Status > 0 {
		whereSql = whereSql + fmt.Sprintf(" and %s.status=:status ", table_name)
	}
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
	updateSql = updateSql + " ,weight=:weight"
	updateSql = updateSql + " ,height=:height"
	updateSql = updateSql + " ,diary=:diary"

	return fmt.Sprintf("update %s set updatetime=:updatetime %s where id=:id and isdelete=false", table_name, updateSql)
}

func delSql() string {
	return fmt.Sprintf("update %s set isdelete=true where id=any(:ids)", table_name)
}

func toggleSql() string {
	return fmt.Sprintf("update %s set disabled=:disabled where id=:id and isdelete=false", table_name)
}
