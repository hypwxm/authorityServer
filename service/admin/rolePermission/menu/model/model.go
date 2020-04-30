package model

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"worldbar/DB/pgsql"
	"worldbar/util/database"
)

type WbAdminRoleMenuPermission struct {
	database.BaseColumns

	RoleId string `json:"roleId" db:"role_id"`
	MenuId string `json:"menuId" db:"menu_id"`
}

type SaveQuery struct {
	RoleId  string `db:"role_id"`
	MenuIds []string
}

func (self *WbAdminRoleMenuPermission) Save(query *SaveQuery) (string, error) {
	var err error
	if strings.TrimSpace(query.RoleId) == "" {
		return "", errors.New(fmt.Sprintf("操作错误"))
	}
	db := pgsql.Open()
	tx, err := db.Beginx()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()
	_, err = tx.Exec(deleteSql(query.RoleId))
	if err != nil {
		return "", err
	}
	for _, v := range query.MenuIds {
		stmt, err := tx.PrepareNamed(saveSql())
		if err != nil {
			return "", err
		}
		var _query = &WbAdminRoleMenuPermission{
			MenuId: v,
			RoleId: query.RoleId,
		}
		_query.BaseColumns.Init()
		log.Println(stmt.QueryString, *_query)

		_, err = stmt.Exec(_query)
		if err != nil {
			return "", err
		}
	}

	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return self.ID, nil
}

type GetQuery struct {
	ID string `db:"id"`
}

type GetModel struct {
	WbAdminRoleMenuPermission
	Like         bool `json:"like" db:"like"`
	TotalLike    int  `json:"totalLike" db:"total_like"`
	TotalComment int  `json:"totalComment" db:"total_comment"`
}

type Query struct {
	pgsql.BaseQuery
	RoleId string `db:"role_id"`
}

type ListModel struct {
	WbAdminRoleMenuPermission
	ParentId string `json:"parentId" db:"parent_id"`
	Name     string `json:"name" db:"name"`
	Path     string `json:"path" db:"path"`
}

func (self *WbAdminRoleMenuPermission) List(query *Query) ([]*ListModel, error) {
	if query == nil {
		query = new(Query)
	}
	db := pgsql.Open()
	fullSql := listSql(query)
	stmt, err := db.PrepareNamed(fullSql)
	if err != nil {
		return nil, err
	}
	log.Println(stmt.QueryString)
	rows, err := stmt.Queryx(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list = make([]*ListModel, 0)
	for rows.Next() {
		var data = new(ListModel)
		err = rows.StructScan(&data)
		if err != nil {
			return nil, err
		}
		list = append(list, data)
	}
	return list, nil
}
