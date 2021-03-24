// 成员和宝宝的关系
package model

import (
	"babygrow/DB/pgsql"
	memberModel "babygrow/service/member/user/model"
	memberService "babygrow/service/member/user/service"

	"babygrow/util/database"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// 用户和宝宝的关系
type GMemberBabyRelation struct {
	database.BaseColumns
	RoleName string `json:"roleName" db:"role_name"`
	BabyId   string `json:"babyId" db:"baby_id"`
	UserId   string `json:"userId" db:"user_id"`
}

func (self *GMemberBabyRelation) Insert(tx *sqlx.Tx) (string, error) {
	var err error

	if strings.TrimSpace(self.RoleName) == "" {
		return "", fmt.Errorf("操作错误")
	}
	if strings.TrimSpace(self.BabyId) == "" {
		return "", fmt.Errorf("操作错误")
	}
	if strings.TrimSpace(self.UserId) == "" {
		return "", fmt.Errorf("操作错误")
	}

	// 插入判断用户登录账号是否已经存在
	stmt, err := tx.PrepareNamed(mbInsertSql())
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

type MbQuery struct {
	pgsql.BaseQuery

	UserId   string `db:"user_id"`
	BabyId   string `db:"baby_id"`
	Keywords string `db:"keywords"`
}

type MbListModel struct {
	GMemberBabyRelation
	Member *memberModel.GMember `json:"member"`
}

func (self *GMemberBabyRelation) List(query *MbQuery) ([]*MbListModel, error) {
	if query == nil {
		query = new(MbQuery)
	}
	db := pgsql.Open()
	fullSql := mbListSql(query)
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

	var list = make([]*MbListModel, 0)
	var ids = make([]string, 0)
	for rows.Next() {
		var item = new(MbListModel)
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
		ids = append(ids, item.UserId)
	}

	members, _, err := memberService.List(&memberModel.Query{
		BaseQuery: pgsql.BaseQuery{
			IDs: ids,
		},
	})

	if err != nil {
		return nil, err
	}

	for _, v := range list {
		for _, vm := range members {
			if v.UserId == vm.ID {
				v.Member = vm
			}
		}
	}
	return list, nil

}

type MBDeleteQuery struct {
	IDs pq.StringArray `db:"ids"`
}

// 删除，批量删除
func (self *GMemberBabyRelation) Delete(query *DeleteQuery) error {
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
	stmt, err := db.PrepareNamed(mbdelSql())
	if err != nil {
		return err
	}
	_, err = stmt.Exec(query)
	return err
}
