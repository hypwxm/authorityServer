package dbModel

import (
	"babygrow/util"

	"gorm.io/gorm"
)

func (u *GDailyComment) TableName() string {
	return "g_member_baby_grow_comment"
}

// 在同一个事务中更新数据
func (u *GDailyComment) AfterDelete(tx *gorm.DB) (err error) {
	// tx.Model(&Address{}).Where("user_id = ?", u.ID).Update("invalid", false)
	return
}

func (u *GDailyComment) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = util.GetUuid()
	return
}
