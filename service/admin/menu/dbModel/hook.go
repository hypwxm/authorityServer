package dbModel

import (
	"authorityServer/util"

	"gorm.io/gorm"
)

func (u *GMenu) TableName() string {
	return BusinessName
}

// 在同一个事务中更新数据
func (u *GMenu) AfterDelete(tx *gorm.DB) (err error) {
	// tx.Model(&Address{}).Where("user_id = ?", u.ID).Update("invalid", false)
	return
}

func (u *GMenu) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = util.GetUuid()
	return
}
