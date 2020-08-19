package model

import (
	"babygrowing/DB/pgsql"
	"fmt"
	"strings"
)

const table_name = "g_my_babies"

func insertSql() string {
	return fmt.Sprintf("insert into %s (createtime, isdelete, disabled, id, name, birthday, gender, avatar, id_card, hobby, good_at, favorite_food, favorite_color, ambition) select :createtime, :isdelete, :disabled, :id, :name, :birthday, :gender, :avatar, :id_card, :hobby, :good_at, :favorite_food, :favorite_color, :ambition returning id", table_name)

}

func listSql(query *Query) (whereSql string, fullSql string) {
	var selectSql = fmt.Sprintf(`SELECT 
				%[1]s.id,
				%[1]s.createtime,
				%[1]s.updatetime,
				%[1]s.name, 
				%[1]s.birthday, 
				%[1]s.gender, 
				%[1]s.avatar,
				%[1]s.id_card, 
				%[1]s.hobby,
				%[1]s.good_at,
				%[1]s.favorite_food, 
				%[1]s.favorite_color, 
				%[1]s.ambition
				FROM %[1]s WHERE 1=1 `, table_name)
	whereSql = pgsql.BaseWhere(query.BaseQuery, table_name)
	if strings.TrimSpace(query.Keywords) != "" {
		// whereSql = whereSql + fmt.Sprintf(" and (%[1]s.title like '%%%[2]s%%' or %[1]s.intro like '%%%[2]s%%' or %[1]s.content like '%%%[2]s%%')", table_name, query.Keywords)
	}

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

func updateSql() string {
	var updateSql = ""
	updateSql = updateSql + " ,name=:name"
	updateSql = updateSql + " ,birthday=:birthday"
	updateSql = updateSql + " ,gender=:gender"
	updateSql = updateSql + " ,avatar=:avatar"
	updateSql = updateSql + " ,id_card=:id_card"
	updateSql = updateSql + " ,hobby=:hobby"
	updateSql = updateSql + " ,good_at=:good_at"
	updateSql = updateSql + " ,favorite_food=:favorite_food"
	updateSql = updateSql + " ,favorite_color=:favorite_color"
	updateSql = updateSql + " ,ambition=:ambition"

	return fmt.Sprintf("update %s set updatetime=:updatetime %s where id=:id and isdelete=false", table_name, updateSql)
}

func delSql() string {
	return fmt.Sprintf("update %s set isdelete=true where id=any(:ids)", table_name)
}

func toggleSql() string {
	return fmt.Sprintf("update %s set disabled=:disabled where id=:id and isdelete=false", table_name)
}
