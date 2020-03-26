package model

import (
	"fmt"
	"strings"
	"worldbar/DB/pgsql"
)

const table_name = "wb_news_dynamics_comment"

func insertSql() string {
	return fmt.Sprintf("insert into %s (createtime, isdelete, disabled, id, title, intro, surface, content, publisher, type) select :createtime, :isdelete, :disabled, :id, :title, :intro, :surface, :content, :publisher, :type returning id", table_name)

}
func listSql(query *Query) (whereSql string, fullSql string) {
	var selectSql = fmt.Sprintf(`SELECT 
				%[0].createtime,
				%[0].updatetime,
				%[0].content,
				%[0].publisher,
				%[0].top_comment_id,
				%[0].prev_comment_id,
				%[0].prev_publisher,
				%[0].news_id,
				%[1].avatar,
				%[1].nickname,
				%[2].avatar as prev_avatar,
				%[2].nickname as prev_nickname
				FROM %[0] left join %[1] on %[0].publisher=%[1].id left join %[2] on %[0].prev_publisher=%[2].id WHERE 1=1 `, table_name, "wb_user", "wb_user")
	whereSql = pgsql.BaseWhere(query.BaseQuery)
	whereSql = whereSql + " and news_id=:news_id"

	if query.PrevCommentId != "" {
		whereSql = whereSql + " and prev_comment_id=:prev_comment_id"
	}
	if query.TopCommentId != "" {
		whereSql = whereSql + " and top_comment_id=:top_comment_id"
	}
	if query.PrevPublisher != "" {
		whereSql = whereSql + " and prev_publisher=:prev_publisher"
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
