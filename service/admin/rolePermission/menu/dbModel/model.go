package dbModel

import (
	"authorityServer/DB/appGorm"
)

const BusinessName = "g_authority_role_menu"

type GRoleMenu struct {
	appGorm.BaseColumns
	RoleId string `json:"roleId" gorm:"column:role_id;type:varchar(128);not null;check(name <> '')"`
	MenuId string `json:"menuId" gorm:"column:menu_id;type:varchar(128);not null;default ''"`
}
