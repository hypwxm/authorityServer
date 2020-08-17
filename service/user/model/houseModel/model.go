package houseModel

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"strings"
	"babygrowing/DB/pgsql"
	"babygrowing/util/database"
)

type WbUserHouse struct {
	database.BaseIDColumns
	EnumsId  string `json:"enumsId" db:"enums_id"`
	OptionId string `json:"optionId" db:"option_id"`
	UserId   string `json:"userId" db:"user_id"`
}

func MultiInsert(list []WbUserHouse, tx *sqlx.Tx) error {
	var err error
	tx, fromOut, err := pgsql.GetOrMakeTx(tx)
	for _, v := range list {
		_, err := v.Insert(tx)
		if err != nil {
			return err
		}
	}
	if !fromOut {
		err = tx.Commit()
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *WbUserHouse) Insert(tx *sqlx.Tx) (string, error) {
	var err error

	if strings.TrimSpace(self.EnumsId) == "" {
		return "", errors.New(fmt.Sprintf("操作错误"))
	}
	if strings.TrimSpace(self.OptionId) == "" {
		return "", errors.New(fmt.Sprintf("操作错误"))
	}
	if strings.TrimSpace(self.UserId) == "" {
		return "", errors.New(fmt.Sprintf("操作错误"))
	}
	// 插入判断用户登录账号是否已经存在
	stmt, err := tx.PrepareNamed(insertSql())
	if err != nil {
		return "", err
	}
	log.Println(stmt.QueryString)
	var lastId string
	self.BaseIDColumns.Init()
	err = stmt.Get(&lastId, self)
	if err != nil {
		return "", err
	}
	return self.ID, nil
}

type GetQuery struct {
	UserId string `db:"user_id"`
}

type Query struct {
	UserId string `json:"userId" db:"user_id"`
}

type ListModel struct {
	WbUserHouse
	OptionName string `json:"optionName" db:"option_name"`
	EnumsName  string `json:"enumsName" db:"enums_name"`
}

func (self *WbUserHouse) GetByUserId(query *GetQuery) ([]*ListModel, error) {
	if query == nil {
		query = new(GetQuery)
	}
	db := pgsql.Open()
	sqlStr := listSql(query)

	stmt, err := db.PrepareNamed(sqlStr)
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
		var item = new(ListModel)
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}

	return list, nil

}
