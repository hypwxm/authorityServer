package model

import (
	"fmt"
	"strings"
	"worldbar/DB/pgsql"
)

const table_name = "wb_news_dynamics"

func insertSql() string {
	return fmt.Sprintf("insert into %s (createtime, isdelete, disabled, id, title, intro, surface, content, publisher, type) select :createtime, :isdelete, :disabled, :id, :title, :intro, :surface, :content, :publisher, :type returning id", table_name)

}

func listSql(query *Query) (whereSql string, fullSql string) {
	var selectSql = fmt.Sprintf(`SELECT 
				%[0].createtime,
				%[0].updatetime,
				%[0].publish_time,
				%[0].title,
				%[0].intro,
				%[0].content,
				%[0].surface,
				%[0].type,
				%[0].sort,
				%[0].status,
				%[0].StatusReason,
				%[0].publisher,
				%[1].avatar,
				%[1].nickname
				FROM %[0] left join %[1] on %[0].publisher=%[1].id WHERE 1=1 `, table_name, "wb_user")
	whereSql = pgsql.BaseWhere(query.BaseQuery)
	if strings.TrimSpace(query.Keywords) != "" {
		whereSql = whereSql + fmt.Sprintf(" and (%[0].title like '%:keywords%' or %[0].intro like '%:keywords%' or %[0].content like '%:keywords%')", table_name)
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
	return fmt.Sprintf("select * from %s where id=:id", table_name)
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
