package model

import (
	"fmt"
	"strings"
	"babygrowing/DB/pgsql"
)

const table_name = "wb_societies"

func insertSql() string {
	return fmt.Sprintf("insert into %s (createtime, isdelete, disabled, id, name, intro, surface, creator, creator_id, people_max, type_id, type_name, status) select :createtime, :isdelete, :disabled, :id, :name, :intro, :surface, :creator, :creator_id, :people_max, :type_id, :type_name, :status returning id", table_name)

}

func listSql(query *Query) (whereSql string, fullSql string) {
	var selectSql = fmt.Sprintf(`SELECT 
				%[1]s.id,
				%[1]s.createtime,
				%[1]s.updatetime,
				%[1]s.publish_time,
				%[1]s.name,
				%[1]s.intro,
				%[1]s.surface,
				%[1]s.creator,
				%[1]s.creator_id,
				%[1]s.people_max,
				%[1]s.type_id,
				%[1]s.type_name,
				%[1]s.status,
				%[1]s.status_reason,
				%[1]s.publish_time
				FROM %[1]s WHERE 1=1 `, table_name)
	whereSql = pgsql.BaseWhere(query.BaseQuery, table_name)
	if strings.TrimSpace(query.Keywords) != "" {
		whereSql = whereSql + fmt.Sprintf(" and (%[1]s.name like '%%%[2]s%%' or %[1]s.intro like '%%%[2]s%%' or %[1]s.creator like '%%%[2]s%%')", table_name, query.Keywords)
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
	updateSql = updateSql + " ,name=:name"
	updateSql = updateSql + " ,intro=:intro"
	updateSql = updateSql + " ,people_max=:people_max"
	updateSql = updateSql + " ,surface=:surface"
	updateSql = updateSql + " ,status=:status"

	return fmt.Sprintf("update %s set updatetime=:updatetime %s where id=:id and isdelete=false", table_name, updateSql)
}

func delSql() string {
	return fmt.Sprintf("update %s set isdelete=true where id=any(:ids)", table_name)
}

func toggleSql() string {
	return fmt.Sprintf("update %s set disabled=:disabled where id=:id and isdelete=false", table_name)
}
