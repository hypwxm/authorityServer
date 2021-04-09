package dbModel

import (
	"babygrow/util"

	"gorm.io/gorm"
)

func (u *Media) TableName() string {
	return "g_media"
}

// 在同一个事务中更新数据
func (u *Media) AfterDelete(tx *gorm.DB) (err error) {
	// tx.Model(&Address{}).Where("user_id = ?", u.ID).Update("invalid", false)
	return
}

func (u *Media) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = util.GetUuid()
	return
}
