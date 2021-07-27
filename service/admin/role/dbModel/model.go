package dbModel

import (
	"github.com/hypwxm/authorityServer/DB/appGorm"
)

const BusinessName = "g_authority_role"

type GRole struct {
	appGorm.BaseColumns
	Name  string `json:"name" gorm:"column:name;type:varchar(20);not null;check(name <> '')"`
	Intro string `json:"intro" gorm:"column:intro;type:varchar(500);not null;default ''"`
	OrgId string `json:"orgId" gorm:"column:org_id;type:varchar(128);not null;check(org_id <> '')"`
}
