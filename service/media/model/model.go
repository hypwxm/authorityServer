package model

import (
	"babygrowing/DB/pgsql"
	"babygrowing/util/database"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Media struct {
	database.BaseColumns
	Url        string `db:"url" json:"url"`
	UserID     string `db:"user_id" json:"userId"`
	Business   string `db:"business" json:"business"`
	BusinessId string `db:"business_id" json:"businessId"`
	Size       int    `db:"size" json:"size"`
}

func (self *Media) Insert(tx ...*sqlx.Tx) (string, error) {
	var err error

	if strings.TrimSpace(self.UserID) == "" {
		return "", errors.New(fmt.Sprintf("操作错误"))
	}
	if strings.TrimSpace(self.Business) == "" {
		return "", errors.New(fmt.Sprintf("未知业务"))
	}
	if strings.TrimSpace(self.BusinessId) == "" {
		return "", errors.New(fmt.Sprintf("未知业务"))
	}

	var ltx *sqlx.Tx
	var stmt *sqlx.NamedStmt
	if len(tx) > 0 {
		ltx = tx[0]
	}
	if ltx != nil {
		stmt, err = ltx.PrepareNamed(insertSql())
	} else {
		db := pgsql.Open()
		stmt, err = db.PrepareNamed(insertSql())
	}

	if err != nil {
		return "", err
	}
	log.Println(stmt.QueryString)
	var lastId string
	self.BaseColumns.Init()
	err = stmt.Get(&lastId, self)
	if err != nil {
		return "", err
	}
	return self.ID, nil
}

func InitMedias(list []*Media, businessName string, businessId string, creator string) []*Media {
	for _, v := range list {
		if v == nil {
			return nil

		}
		v.Init()
		v.Business = businessName
		v.BusinessId = businessId
		// 有一种情况是之前人创建的虽然重新插入，但是前端把之前的创建者传过来了，就重新保存下，不能覆盖
		if v.UserID == "" {
			v.UserID = creator
		}
	}
	return list
}

func StoreMedias(list []*Media, tx ...*sqlx.Tx) error {
	var err error
	for _, v := range list {
		if v == nil {
			return fmt.Errorf("文件错误")
		}
		_, err = v.Insert(tx...)
		if err != nil {
			return err
		}
	}
	return nil
}

type Query struct {
	pgsql.BaseQuery
	BusinessIds pq.StringArray `json:"businessIds" db:"business_ids"`
	Businesses  pq.StringArray `json:"businesses" db:"businesses"`
}

func (self *Media) List(query *Query) ([]*Media, int, error) {
	if query == nil {
		query = new(Query)
	}
	db := pgsql.Open()
	whereSql, fullSql := listSql(query)
	// 以上部分为查询条件，接下来是分页和排序
	count, err := self.GetCount(db, query, whereSql)
	if err != nil {
		return nil, 0, err
	}
	stmt, err := db.PrepareNamed(fullSql)
	if err != nil {
		return nil, 0, err
	}
	log.Println(stmt.QueryString)

	rows, err := stmt.Queryx(query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list = make([]*Media, 0)
	for rows.Next() {
		var item = new(Media)
		err = rows.StructScan(&item)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, item)
	}

	return list, count, nil

}

func (self *Media) GetCount(db *sqlx.DB, query *Query, whereSql ...string) (int, error) {
	if query == nil {
		query = new(Query)
	}
	sqlStr := countSql(whereSql...)
	stmt, err := db.PrepareNamed(sqlStr)
	if err != nil {
		return 0, err
	}
	var count int
	err = stmt.Get(&count, query)
	log.Println(stmt.QueryString, query)
	return count, err
}

type DeleteQuery struct {
	IDs         pq.StringArray `db:"ids"`
	Businesses  pq.StringArray `json:"businesses" db:"businesses"`
	BusinessIds pq.StringArray `json:"businessIds" db:"business_ids"`
}

// 删除，批量删除
func (self *Media) Delete(query *DeleteQuery, tx ...*sqlx.Tx) error {
	if query == nil {
		return errors.New("无操作条件")
	}
	if len(query.IDs) == 0 && len(query.BusinessIds) == 0 {
		return errors.New("操作条件错误")
	}

	var ltx *sqlx.Tx
	if len(tx) > 0 {
		ltx = tx[0]
	}
	var stmt *sqlx.NamedStmt
	var err error
	if ltx != nil {
		// 通过事务来删除
		stmt, err = ltx.PrepareNamed(delSql(query))
	} else {
		db := pgsql.Open()
		stmt, err = db.PrepareNamed(delSql(query))
	}
	if err != nil {
		return err
	}
	_, err = stmt.Exec(query)
	return err
}