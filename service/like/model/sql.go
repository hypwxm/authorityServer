package model

import (
	"fmt"
	"strings"
	"babygrowing/DB/pgsql"
)

const table_name = "wb_like"

func insertSql(query *WbLike) string {
	return fmt.Sprintf(`insert into %s 
								(createtime, isdelete, disabled, id, source_type, source_id) 
								select :createtime, :isdelete, :disabled, :id, :user_id, :source_type, :source_id 
								where not exists(select 1 from %s where source_type=%d and source_id=%s and isdelete='false') returning id`,
		table_name, table_name, query.SourceType, query.SourceId)

}

func listSql(query *Query) (whereSql string, fullSql string) {
	var selectSql = ""
	whereSql = pgsql.BaseWhere(query.BaseQuery)

	if query.SourceType == SourceTypeNews {
		selectSql = fmt.Sprintf(`SELECT 
				%[1]s.createtime,
				%[1]s.user_id,
				%[1]s.source_type,
				%[1]s.source_id,
				%[2]s.avatar,
				%[2]s.nickname,
				%[3]s.title as news_title,
				%[3]s.surface as news_surface
				FROM %[1]s left join %[2]s on %[1]s.user_id=%[2]s.id left join %[3]s on %[3]s.id=%[1]s.source_id WHERE 1=1 `, table_name, "wb_user", "wb_news_dynamics")
		whereSql = whereSql + fmt.Sprintf(" and %[1]s.source_type=%[2]d", table_name, SourceTypeNews)
	} else if query.SourceType == SourceTypeUser {
		selectSql = fmt.Sprintf(`SELECT 
				%[1]s.createtime,
				%[1]s.user_id,
				%[1]s.source_type,
				%[1]s.source_id,
				%[2]s.avatar,
				%[2]s.nickname,
				%[3]s.avatar as like_avatar,
				%[3]s.nickname as like_nickname,
				FROM %[1]s left join %[2]s on %[1]s.user_id=%[2]s.id left join %[3]s on %[3]s.id=%[1]s.source_id WHERE 1=1 `, table_name, "wb_user", "wb_user")
		whereSql = whereSql + fmt.Sprintf(" and %[1]s.source_type=%[2]d", table_name, SourceTypeUser)
	} else if query.SourceType == SourceTypeComment {
		selectSql = fmt.Sprintf(`SELECT 
				%[1]s.createtime,
				%[1]s.user_id,
				%[1]s.source_type,
				%[1]s.source_id,
				%[2]s.avatar,
				%[2]s.nickname,
				%[3]s.content as commentContent,
				FROM %[1]s left join %[2]s on %[1]s.user_id=%[2]s.id left join %[3]s on %[3]s.id=%[1]s.source_id WHERE 1=1 `, table_name, "wb_user", "wb_news_dynamics_comment")
		whereSql = whereSql + fmt.Sprintf(" and %[1]s.source_type=%[2]d", table_name, SourceTypeComment)
	} else {
		selectSql = fmt.Sprintf(`SELECT 
				%[1]s.createtime,
				%[1]s.user_id,
				%[1]s.source_type,
				%[1]s.source_id,
				%[2]s.avatar,
				%[2]s.nickname,
				%[3]s.title as news_title,
				%[3]s.surface as news_surface,
				%[4]s.avatar as like_avatar,
				%[4]s.surface as like_nickname
				FROM %[1]s left join %[2]s on %[1]s.user_id=%[2]s.id left join %[3]s on %[3]s.id=%[1]s.source_id and %[1]s.source_type=1 left join %[4]s on %[4]s.id=%[1]s.source_id and %[1]s.source_type=2 WHERE 1=1 `, table_name, "wb_user", "wb_news_dynamics", "wb_user")
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
