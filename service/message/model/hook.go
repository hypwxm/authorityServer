package model

import (
	"babygrow/util"

	"gorm.io/gorm"
)

// 在同一个事务中更新数据
func (u *GMessage) AfterDelete(tx *gorm.DB) (err error) {
	// tx.Model(&Address{}).Where("user_id = ?", u.ID).Update("invalid", false)
	return
}

func (u *GMessage) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = util.GetUuid()
	return
}
