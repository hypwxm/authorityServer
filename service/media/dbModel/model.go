package dbModel

import (
	"authorityServer/DB/appGorm"
)

const BusinessName = "g_media"

type Media struct {
	appGorm.BaseColumns
	Url        string `db:"url" json:"url" gorm:"column:url;type:varchar(500);not null;default ''"`
	UserID     string `db:"user_id" json:"userId" gorm:"column:user_id;type:varchar(128);not null;check(user_id <> '')"`
	Business   string `db:"business" json:"business" gorm:"column:business;type:varchar(50);not null;default ''"`
	BusinessId string `db:"business_id" json:"businessId" gorm:"column:business_id;type:varchar(128);not null;check(business_id <> '')"`
	Size       int    `db:"size" json:"size" gorm:"column:size;not null;default 0"`
}
