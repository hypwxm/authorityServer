package model

import (
	"fmt"
	"strings"
	"babygrowing/DB/pgsql"
	"babygrowing/service/like/model"
)

const table_name = "wb_news_dynamics_comment"

func insertSql() string {
	return fmt.Sprintf("insert into %s (createtime, isdelete, disabled, id, title, intro, surface, content, publisher, type) select :createtime, :isdelete, :disabled, :id, :title, :intro, :surface, :content, :publisher, :type returning id", table_name)

}
func listSql(query *Query) (whereSql string, fullSql string) {
	var selectSql = fmt.Sprintf(`SELECT 
				%[1]s.createtime,
				%[1]s.updatetime,
				%[1]s.content,
				%[1]s.publisher,
				%[1]s.top_comment_id,
				%[1]s.prev_comment_id,
				%[1]s.prev_publisher,
				%[1]s.news_id,
				%[2]s.avatar,
				%[2]s.nickname,
				%[3]s.avatar as prev_avatar,
				%[3]s.nickname as prev_nickname,
				case when %[4]s.id <> null then true else false end as like
				FROM %[1]s left join %[2]s on %[1]s.publisher=%[2]s.id left join %[3]s on %[1]s.prev_publisher=%[3]s.id left join %[4]s on %[4]s.source_id=%[1]s.id and %[4]s.source_type=%[5]d WHERE 1=1 `, table_name, "wb_user", "wb_user", "wb_like", model.SourceTypeComment)
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
