// 成员和宝宝的关系
package model

import (
	"babygrow/DB/appGorm"
	"babygrow/DB/pgsql"
	memberModel "babygrow/service/member/user/model"
	memberService "babygrow/service/member/user/service"

	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// 用户和宝宝的关系
type GMemberBabyRelation struct {
	appGorm.BaseColumns
	RoleName string `json:"roleName" db:"role_name" gorm:"column:role_name;type:varchar(10);not null;check(role_name <> '')"`
	BabyId   string `json:"babyId" db:"baby_id" gorm:"column:baby_id;type:varchar(128);not null;check(baby_id <> '');uniqueIndex:user_baby_id"`
	UserId   string `json:"userId" db:"user_id" gorm:"column:user_id;type:varchar(128);not null;check(user_id <> '');uniqueIndex:user_baby_id"`

	Account string `json:"account" db:"-" gorm:"-"`
}

func (self *GMemberBabyRelation) Insert(tx *gorm.DB) (string, error) {
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

	localTx := false
	if tx == nil {
		localTx = true
		db := appGorm.Open()
		tx = db.Begin()
		if err := tx.Error; err != nil {
			return "", err
		}
	}

	err = tx.Model(&GMemberBabyRelation{}).Create(self).Error
	if err != nil {
		return "", err
	}
	if localTx {
		err = tx.Commit().Error
		if err != nil {
			return "", err
		}
	}

	return self.ID, nil
}

type MbQuery struct {
	appGorm.BaseQuery

	UserId   string `db:"user_id"`
	BabyId   string `db:"baby_id"`
	Keywords string `db:"keywords"`
}

type MbListModel struct {
	GMemberBabyRelation
	Member *memberModel.GMember `json:"member" gorm:"-"`
}

func (self *GMemberBabyRelation) List(query *MbQuery) ([]*MbListModel, error) {
	if query == nil {
		query = new(MbQuery)
	}
	db := appGorm.Open()
	tx := db.Table("g_member_baby_relation").Select(`*`)
	tx.Scopes(appGorm.BaseWhere(query.BaseQuery))
	if query.UserId != "" {
		tx.Where("user_id=?", query.UserId)
	}

	if query.BabyId != "" {
		tx.Where("baby_id=?", query.BabyId)
	}
	var list = make([]*MbListModel, 0)
	err := tx.Scopes(appGorm.Paginate(query.BaseQuery)).Find(&list).Error
	if err != nil {
		return nil, err
	}
	var ids = make([]string, len(list))
	for _, v := range list {
		ids = append(ids, v.UserId)
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
	IDs pq.StringArray `json:"ids" db:"ids"`
}

// 删除，批量删除
func (self *GMemberBabyRelation) Delete(query *MBDeleteQuery) error {
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
	db := appGorm.Open()
	return db.Where("id=any(?)", query.IDs).Delete(&GMemberBabyRelation{}).Unscoped().Error
}
