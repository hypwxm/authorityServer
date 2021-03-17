package model

import (
	"fmt"
	"io/ioutil"
	"strings"
)

const table_name = "g_member_family"

func GetSqlFile() ([]byte, error) {
	b, err := ioutil.ReadFile("scheme.sql")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func insertSql() string {
	return fmt.Sprintf("insert into %s (createtime, isdelete, disabled, id, name, creator, label, intro) select :createtime, :isdelete, :disabled, :id, :name, :creator, :label, :intro returning id", table_name)

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

func updateSql() string {
	var updateSql = ""
	updateSql = updateSql + " ,name=:name"

	return fmt.Sprintf("update %s set updatetime=:updatetime %s where id=:id and isdelete=false", table_name, updateSql)
}

func delSql() string {
	return fmt.Sprintf("update %s set isdelete=true where id=any(:ids)", table_name)
}

func toggleSql() string {
	return fmt.Sprintf("update %s set disabled=:disabled where id=:id and isdelete=false", table_name)
}
