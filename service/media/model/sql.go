package model

import (
	"babygrowing/DB/pgsql"
	"fmt"
	"io/ioutil"
	"strings"
)

const table_name = "media"

func GetSqlFile() ([]byte, error) {
	b, err := ioutil.ReadFile("scheme.sql")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func insertSql() string {
	return fmt.Sprintf("insert into %s (createtime, isdelete, disabled, id, url, user_id, business) select :createtime, :isdelete, :disabled, :id, :url, :user_id, :business returning id", table_name)

}

func listSql(query *Query) (whereSql string, fullSql string) {
	var selectSql = fmt.Sprintf(`SELECT 
				%[1]s.id,
				%[1]s.createtime,
				%[1]s.updatetime,
				%[1]s.weight,
				%[1]s.height,
				%[1]s.diary,
				%[1]s.user_id,
				%[1]s.baby_id
				FROM %[1]s WHERE 1=1 `, table_name)
	whereSql = pgsql.BaseWhere(query.BaseQuery, table_name)

	if query.OrderBy == "" {
		// query.OrderBy = "sort asc"
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

func delSql() string {
	return fmt.Sprintf("update %s set isdelete=true where id=any(:ids)", table_name)
}
