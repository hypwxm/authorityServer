package dbModel

import (
	"babygrow/util"

	"gorm.io/gorm"
)

func (u *GMember) TableName() string {
	return "g_member"
}

// 在同一个事务中更新数据
func (u *GMember) AfterDelete(tx *gorm.DB) (err error) {
	// tx.Model(&Address{}).Where("user_id = ?", u.ID).Update("invalid", false)
	return
}

func (u *GMember) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = util.GetUuid()
	return
}
