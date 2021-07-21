package dbModel

import (
	"authorityServer/DB/appGorm"
)

const BusinessName = "g_authority_org"

type GOrg struct {
	appGorm.BaseColumns
	Name     string `json:"name" db:"name" gorm:"column:name;type:varchar(20);not null;check(name <> '')"`
	Sort     int    `json:"sort" db:"sort"`
	ParentId string `json:"parentId" gorm:"column:parent_id;type:varchar(128);not null;default ''"`
}
