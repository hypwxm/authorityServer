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
	Url      string `db:"url" json:"url"`
	UserID   string `db:"user_id" json:"userId"`
	Business string `db:"business" json:"business"`
}

func (self *Media) Insert() (string, error) {
	var err error

	if strings.TrimSpace(self.UserID) == "" {
		return "", errors.New(fmt.Sprintf("操作错误"))
	}
	if strings.TrimSpace(self.Business) == "" {
		return "", errors.New(fmt.Sprintf("未知业务"))
	}

	db := pgsql.Open()

	stmt, err := db.PrepareNamed(insertSql())
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

func InitMedias(list []*Media, businessName string, creator string) []*Media {
	for _, v := range list {
		v.Init()
		v.Business = businessName
		v.UserID = creator
	}
	return list
}

func StoreMedias(list []*Media) error {
	var err error
	for _, v := range list {
		_, err = v.Insert()
		if err != nil {
			return err
		}
	}
	return nil
}

type Query struct {
	pgsql.BaseQuery
}

type ListModel struct {
	Media
}

func (self *Media) List(query *Query) ([]*ListModel, int64, error) {
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

	var list = make([]*ListModel, 0)
	for rows.Next() {
		var item = new(ListModel)
		err = rows.StructScan(&item)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, item)
	}

	return list, count, nil

}

func (self *Media) GetCount(db *sqlx.DB, query *Query, whereSql ...string) (int64, error) {
	if query == nil {
		query = new(Query)
	}
	sqlStr := countSql(whereSql...)
	stmt, err := db.PrepareNamed(sqlStr)
	if err != nil {
		return 0, err
	}
	var count int64
	err = stmt.Get(&count, query)
	log.Println(stmt.QueryString, query)
	return count, err
}

type DeleteQuery struct {
	IDs pq.StringArray `db:"ids"`
}

// 删除，批量删除
func (self *Media) Delete(query *DeleteQuery) error {
	if query == nil {
		return errors.New("无操作条件")
	}
	if len(query.IDs) == 0 {
		return errors.New("操作条件错误")
	}
	for _, v := range query.IDs {
		if strings.TrimSpace(v) == "" {
			return errors.New("操作条件错误")
		}
	}

	db := pgsql.Open()
	stmt, err := db.PrepareNamed(delSql())
	if err != nil {
		return err
	}
	_, err = stmt.Exec(query)
	return err
}
