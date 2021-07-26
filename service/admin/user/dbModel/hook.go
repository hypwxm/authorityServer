package dbModel

import (
	"github.com/hypwxm/authorityServer/util"

	"gorm.io/gorm"
)

func (u *GAdminUser) TableName() string {
	return BusinessName
}

// 在同一个事务中更新数据
func (u *GAdminUser) AfterDelete(tx *gorm.DB) (err error) {
	// tx.Model(&Address{}).Where("user_id = ?", u.ID).Update("invalid", false)
	return
}

func (u *GAdminUser) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = util.GetUuid()
	return
}

func (u *GUserRole) TableName() string {
	return "g_authority_user_role"
}

// 在同一个事务中更新数据
func (u *GUserRole) AfterDelete(tx *gorm.DB) (err error) {
	// tx.Model(&Address{}).Where("user_id = ?", u.ID).Update("invalid", false)
	return
}

func (u *GUserRole) BeforeCreate(tx *gorm.DB) (err error) {
	return
}
