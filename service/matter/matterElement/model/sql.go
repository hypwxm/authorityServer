package model

import (
	"fmt"
	"strings"
	"worldbar/DB/pgsql"
)

const table_name = "wb_matter_element"

func insertSql() string {
	return fmt.Sprintf(`insert into %s 
	(createtime, isdelete, disabled, id, title, intro, min, max, type, matter_id) 
	select :createtime, :isdelete, :disabled, :id, :title, :intro, :min, :max, :type, :matter_id returning id`, table_name)

}

func listSql(query *Query) (whereSql string, fullSql string) {
	var selectSql = fmt.Sprintf(`SELECT 
				%[1]s.id,
				%[1]s.createtime,
				%[1]s.title,
				%[1]s.intro,
				%[1]s.type,
				%[1]s.matter_id,
				%[1]s.min,
				%[1]s.max,
				%[1]s.disabled
				FROM %[1]s WHERE 1=1 `, table_name)
	whereSql = pgsql.BaseWhere(query.BaseQuery)
	if strings.TrimSpace(query.Keywords) != "" {
		whereSql = whereSql + fmt.Sprintf(" and (%[1]s.title like '%%:keywords%%' or %[1]s.intro like '%%:keywords%%')", table_name)
	}

	query.OrderBy = "createtime asc"
	optionSql := pgsql.BaseOption(query.BaseQuery)
	return whereSql, selectSql + whereSql + optionSql
}

func countSql(whereSql ...string) string {
	return fmt.Sprintf("select count(*) from %s where 1=1 %s", table_name, strings.Join(whereSql, " "))
}

func getByIdSql() string {
	return fmt.Sprintf(`
			select 
				%[1]s.*,
				where id=:id and isdelete=false`,
		table_name)
}

func updateSql() string {
	var updateSql = ""
	updateSql = updateSql + " ,title=:title"
	updateSql = updateSql + " ,intro=:intro"
	updateSql = updateSql + " ,min=:min"
	updateSql = updateSql + " ,max=:max"

	return fmt.Sprintf("update %s set updatetime=:updatetime %s where id=:id and isdelete=false", table_name, updateSql)
}

func delSql() string {
	return fmt.Sprintf("update %s set isdelete=true where id=any(:ids)", table_name)
}

func toggleSql() string {
	return fmt.Sprintf("update %s set disabled=:disabled where id=:id and isdelete=false", table_name)
}
