package model

import (
	"fmt"
	"strings"
	"babygrowing/DB/pgsql"
	"babygrowing/service/like/model"
)

const table_name = "wb_matter_element_option"

func insertSql() string {
	return fmt.Sprintf("insert into %s (createtime, isdelete, disabled, id, title, matter_id, element_id) select :createtime, :isdelete, :disabled, :id, :title, :matter_id, :element_id returning id", table_name)
}

func listSql(query *Query) (whereSql string, fullSql string) {
	var selectSql = fmt.Sprintf(`SELECT 
				%[1]s.id,
				%[1]s.createtime,
				%[1]s.title,
				%[1]s.matter_id,
				%[1]s.element_id
				FROM %[1]s WHERE 1=1 `, table_name)
	whereSql = pgsql.BaseWhere(query.BaseQuery)
	if strings.TrimSpace(query.Keywords) != "" {
		whereSql = whereSql + fmt.Sprintf(" and (%[1]s.title like '%%:keywords%%' or %[1]s.intro like '%%:keywords%%' or %[1]s.content like '%%:keywords%%')", table_name)
	}

	if query.MatterId != "" {
		whereSql = whereSql + fmt.Sprintf(" and %s.matter_id=:matter_id", table_name)
	}
	if query.ElementId != "" {
		whereSql = whereSql + fmt.Sprintf(" and %s.element_id=:element_id", table_name)
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
				(select count(*) from %[2]s where source_type=%[3]d and source_id=:id) as total_like,
				(select count(*) from %[4]s where news_id=:id and isdelete=false) as total_comment, 
				case when %[2]s.id <> null then true else false end as like 
				from %[1]s left join %[2]s on %[1]s.id=%[2]s.source_id and %[2]s.source_type=%[3]d 
				where id=:id and isdelete=false`,
		table_name, "wb_like", model.SourceTypeNews, "wb_news_dynamics_comment")
}

func updateSql() string {
	var updateSql = ""
	updateSql = updateSql + " ,title=:title"
	return fmt.Sprintf("update %s set updatetime=:updatetime %s where id=:id and isdelete=false", table_name, updateSql)
}

func delSql(query *DeleteQuery) string {
	if query.ElementId != "" && len(query.IDs) == 0 {
		return fmt.Sprintf("update %s set isdelete=true where element_id=:element_id", table_name)
	} else if query.ElementId != "" && len(query.IDs) > 0 {
		return fmt.Sprintf("update %s set isdelete=true where element_id=:element_id and id=any(:ids)", table_name)
	} else {
		return fmt.Sprintf("update %s set isdelete=true where id=any(:ids)", table_name)
	}
}

func toggleSql() string {
	return fmt.Sprintf("update %s set disabled=:disabled where id=:id and isdelete=false", table_name)
}
