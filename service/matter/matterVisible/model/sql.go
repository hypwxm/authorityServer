package model

import (
	"fmt"
	"babygrowing/service/like/model"
)

const table_name = "wb_matter_visible"

func insertSql() string {
	return fmt.Sprintf("insert into %s (matter_id, user_id) select :matter_id, :user_id", table_name)

}

func listSql() (whereSql string) {
	var selectSql = fmt.Sprintf(`SELECT 
				%[1]s.matter_id,
				%[1]s.user_id
				FROM %[1]s WHERE matter_id=:matter_id `, table_name)
	return selectSql
}

func countSql() string {
	return fmt.Sprintf("select count(*) from %s where matter_id=:matter_id", table_name)
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
	return fmt.Sprintf("delete from %s where user_id=:user_id and matter_id=:matter_id", table_name)
}

func toggleSql() string {
	return fmt.Sprintf("update %s set disabled=:disabled where id=:id and isdelete=false", table_name)
}
