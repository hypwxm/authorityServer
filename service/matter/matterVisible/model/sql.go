package model

import (
	"fmt"
	"strings"
	"worldbar/DB/pgsql"
	"worldbar/service/like/model"
)

const table_name = "wb_matter_visible"

func insertSql() string {
	return fmt.Sprintf("insert into %s (matter_id, user_id) select :matter_id, :user_id", table_name)

}

func listSql(query *Query) (whereSql string, fullSql string) {
	var selectSql = fmt.Sprintf(`SELECT 
				%[1]s.createtime,
				%[1]s.updatetime,
				%[1]s.publish_time,
				%[1]s.title,
				%[1]s.intro,
				%[1]s.content,
				%[1]s.surface,
				%[1]s.type,
				%[1]s.sort,
				%[1]s.status,
				%[1]s.StatusReason,
				%[1]s.publisher,
				%[2]s.avatar,
				%[2]s.nickname,
				case when %[3]s.id <> null then true else false end as like
				FROM %[1]s left join %[2]s on %[1]s.publisher=%[2]s.id left join %[3]s on %[3]s.source_id=%[1]s.id and %[3]s.source_type=1 WHERE 1=1 `, table_name, "wb_user", "wb_news_dynamics_comment")
	whereSql = pgsql.BaseWhere(query.BaseQuery)
	if strings.TrimSpace(query.Keywords) != "" {
		whereSql = whereSql + fmt.Sprintf(" and (%[1]s.title like '%%:keywords%%' or %[1]s.intro like '%%:keywords%%' or %[1]s.content like '%%:keywords%%')", table_name)
	}

	if query.Status > 0 {
		whereSql = whereSql + " and status=:status "
	}
	if query.OrderBy == "" {
		query.OrderBy = "sort asc"
	}
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
	updateSql = updateSql + " ,intro=:intro"
	updateSql = updateSql + " ,content=:content"
	updateSql = updateSql + " ,surface=:surface"

	return fmt.Sprintf("update %s set updatetime=:updatetime %s where id=:id and isdelete=false", table_name, updateSql)
}

func delSql() string {
	return fmt.Sprintf("update %s set isdelete=true where id=any(:ids)", table_name)
}

func toggleSql() string {
	return fmt.Sprintf("update %s set disabled=:disabled where id=:id and isdelete=false", table_name)
}