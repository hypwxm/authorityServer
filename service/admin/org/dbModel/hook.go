package dbModel

import (
	"github.com/hypwxm/authorityServer/util"

	"gorm.io/gorm"
)

func (u *GOrg) TableName() string {
	return BusinessName
}

// 在同一个事务中更新数据
func (u *GOrg) AfterDelete(tx *gorm.DB) (err error) {
	// tx.Model(&Address{}).Where("user_id = ?", u.ID).Update("invalid", false)
	return
}

func (u *GOrg) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = util.GetUuid()
	return
}
