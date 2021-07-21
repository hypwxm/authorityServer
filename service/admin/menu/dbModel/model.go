package dbModel

import (
	"authorityServer/DB/appGorm"
)

const BusinessName = "g_authority_menu"

type GMenu struct {
	appGorm.BaseColumns

	Name string `json:"name" db:"name" gorm:"column:name;type:varchar(20);not null;check(name <> '')"`
	Icon string `json:"icon" db:"icon" gorm:"column:icon;type:varchar(20);not null;default ''"`
	Path string `json:"path" db:"path" gorm:"column:path;type:varchar(100);not null;check(path <> '')"`

	Sort     int    `json:"sort" db:"sort"`
	ParentId string `json:"parentId" db:"parent_id" gorm:"column:parent_id;type:varchar(128);not null;default ''"`

	UserId string `json:"userId" db:"-" gorm:"-"`
}
