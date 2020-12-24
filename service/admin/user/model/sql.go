package model

import (
	"babygrowing/DB/pgsql"
	"fmt"
	"strings"
)

const table_name = "g_admin_user"

func insertSql() string {
	return fmt.Sprintf(`insert into %s 
	(createtime, isdelete, disabled, id, account, password, username, salt, avatar, post)
	select :createtime, :isdelete, :disabled, :id, :account, :password, :username, :salt, :type, :avatar, :post
	where not exists(select 1 from %[1]s where account=:account and isdelete=false) returning id`,
		table_name,
	)
}

func listSql(query *Query) (whereSql string, fullSql string) {
	var selectSql = fmt.Sprintf(`SELECT 
				%[1]s.id,
				%[1]s.createtime,
				%[1]s.updatetime,
				%[1]s.account,
				%[1]s.username,
				%[1]s.avatar,
				%[1]s.post,
				%[1]s.disabled,
				%[1]s.sort
				FROM %[1]s WHERE 1=1 and %[1]s.isdelete=false`, table_name)
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
				%[1]s.id,
				%[1]s.createtime,
				%[1]s.updatetime,
				%[1]s.account,
				%[1]s.username,
				%[1]s.avatar,
				%[1]s.post,
				%[1]s.disabled,
				%[1]s.sort
				from %[1]s
				where %[1]s.id=:id and %[1]s.isdelete=false`,
		table_name, "wb_admin_role")
}

func updateSql(query *UpdateByIDQuery) string {
	var updateSql = ""
	updateSql = updateSql + " ,username=:username"
	updateSql = updateSql + " ,avatar=:avatar"
	updateSql = updateSql + " ,post=:post"
	updateSql = updateSql + " ,sort=:sort"
	updateSql = updateSql + " ,contact_way=:contact_way"

	if query.Password != "" {
		updateSql = updateSql + " ,password=:password"
	}

	return fmt.Sprintf("update %s set updatetime=:updatetime %s where id=:id and isdelete=false", table_name, updateSql)
}

func delSql() string {
	return fmt.Sprintf("update %s set isdelete=true where id=any(:ids)", table_name)
}

func toggleSql() string {
	return fmt.Sprintf("update %s set disabled=:disabled where id=:id and isdelete=false", table_name)
}
