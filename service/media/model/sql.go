package model

import (
	"fmt"
	"strings"

	"github.com/hypwxm/authorityServer/DB/pgsql"
)

// 表名
const table_name = "g_media"

func insertSql() string {
	return fmt.Sprintf("insert into %s (createtime, isdelete, disabled, id, url, user_id, business, business_id, size) select :createtime, :isdelete, :disabled, :id, :url, :user_id, :business, :business_id, :size returning id", table_name)

}

func listSql(query *Query) (whereSql string, fullSql string) {
	var selectSql = fmt.Sprintf(`SELECT 
				%[1]s.*
				FROM %[1]s WHERE 1=1 `, table_name)
	whereSql = pgsql.BaseWhere(query.BaseQuery, table_name)

	if query.OrderBy == "" {
		// query.OrderBy = "sort asc"
	}

	if len(query.Businesses) > 0 {
		whereSql = whereSql + fmt.Sprintf(" and business=any(:businesses)")
	}
	if len(query.BusinessIds) > 0 {
		whereSql = whereSql + fmt.Sprintf(" and business_id=any(:business_ids)")
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

func delSql(query *DeleteQuery) string {
	var whereSQL = " 1=1 "
	if len(query.Businesses) > 0 {
		whereSQL = whereSQL + fmt.Sprintf(" and business=any(:businesses)")
	}
	if len(query.BusinessIds) > 0 {
		whereSQL = whereSQL + fmt.Sprintf(" and business_id=any(:business_ids)")
	}
	if len(query.IDs) > 0 {
		whereSQL = whereSQL + fmt.Sprintf(" and id=any(:ids)")
	}
	return fmt.Sprintf("update %s set isdelete=true where %s", table_name, whereSQL)
}
