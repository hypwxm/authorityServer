package model

import (
	"babygrow/util"

	"gorm.io/gorm"
)

func (u *GMemberBabyRelation) TableName() string {
	return "g_member_baby_relation"
}

// 在同一个事务中更新数据
func (u *GMemberBabyRelation) AfterDelete(tx *gorm.DB) (err error) {
	// tx.Model(&Address{}).Where("user_id = ?", u.ID).Update("invalid", false)
	return
}

func (u *GMemberBabyRelation) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = util.GetUuid()
	return
}