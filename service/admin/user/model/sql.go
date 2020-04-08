package model

import (
	"fmt"
	"strings"
	"worldbar/DB/pgsql"
	"worldbar/service/like/model"
)

const table_name = "wb_admin_user"

func insertSql() string {
	return fmt.Sprintf(`insert into %s 
	(createtime, isdelete, disabled, id, account, password, username, salt, type, avatar, role_id)
	select :createtime, :isdelete, :disabled, :id, :account, :password, :username, :salt, :type, :avatar, :role_id returning id`,
		table_name,
	)
}

func listSql(query *Query) (whereSql string, fullSql string) {
	var selectSql = fmt.Sprintf(`SELECT 
				%[1]s.createtime,
				%[1]s.updatetime,
				%[1]s.account,
				%[1]s.username,
				%[1]s.avatar,
				%[1]s.type,
				%[1]s.disabled,
				%[2]s.id as role_id,
				%[2]s.name as role_name
				FROM %[1]s left join %[2]s on %[1]s.role_id=%[2]s.id WHERE 1=1 `, table_name, "wb_admin_role")
	whereSql = pgsql.BaseWhere(query.BaseQuery, table_name)
	if strings.TrimSpace(query.Keywords) != "" {
		whereSql = whereSql + fmt.Sprintf(" and (%[1]s.account like '%%:keywords%%' or %[1]s.username like '%%:keywords%%')", table_name)
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
