package model

import (
	"fmt"
	"io/ioutil"
	"strings"
)

const table_name = "g_member_baby"
const table_name_mb = "g_member_baby_relation"

func GetSqlFile() ([]byte, error) {
	b, err := ioutil.ReadFile("scheme.sql")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func insertSql() string {
	return fmt.Sprintf("insert into %s (createtime, isdelete, disabled, id, name, birthday, gender, avatar, id_card, hobby, good_at, favorite_food, favorite_color, ambition, user_id, weight, height) select :createtime, :isdelete, :disabled, :id, :name, :birthday, :gender, :avatar, :id_card, :hobby, :good_at, :favorite_food, :favorite_color, :ambition, :user_id, :weight, :height returning id", table_name)
}

func mbInsertSql() string {
	return fmt.Sprintf("insert into %s (createtime, isdelete, disabled, id, role_name, baby_id, user_id) select :createtime, :isdelete, :disabled, :id, :role_name, :baby_id, :user_id returning id", table_name_mb)
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
	updateSql = updateSql + " ,weight=:weight"
	updateSql = updateSql + " ,height=:height"

	return fmt.Sprintf("update %s set updatetime=:updatetime %s where id=:id and isdelete=false", table_name, updateSql)
}

func delSql() string {
	return fmt.Sprintf("update %s set isdelete=true where id=any(:ids)", table_name)
}

func toggleSql() string {
	return fmt.Sprintf("update %s set disabled=:disabled where id=:id and isdelete=false", table_name)
}

func mbdelSql() string {
	return fmt.Sprintf("update %s set isdelete=true where id=any(:ids)", table_name_mb)
}
