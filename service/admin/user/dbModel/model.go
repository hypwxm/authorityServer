package dbModel

import (
	"github.com/hypwxm/authorityServer/DB/appGorm"
)

const BusinessName = "g_authority_user"

type GAdminUser struct {
	appGorm.BaseColumns
	CreatorId string `json:"creatorId" gorm:"column:creator_id;type:varchar(128);not null;check(creator_id <> '')"`
	Creator   string `json:"creator" gorm:"column:creator;type:varchar(100);not null;default ''"`

	Account    string `json:"account" gorm:"column:account;type:varchar(100);not null;check(account <> '');uniqueIndex"`
	Password   string `json:"password" gorm:"column:password;type:varchar(250);not null;check(password <> '')"`
	Username   string `json:"username" gorm:"column:username;type:varchar(100);not null;check(username <> '')"`
	ContactWay string `json:"contactWay" gorm:"column:contact_way;type:varchar(100);not null;default ''"`
	Post       string `json:"post" gorm:"column:post;type:varchar(100);not null;default ''"`
	Salt       string `json:"salt" gorm:"column:salt;type:varchar(100);not null;check(salt <> '')"`
	Sort       int    `json:"sort" gorm:"column:sort;type:int;not null;default 0"`
}

type GUserRole struct {
	UserId string `json:"userId" gorm:"column:user_id;type:varchar(128);not null;check(user_id <> '')"`
	RoleId string `json:"roleId" gorm:"column:role_id;type:varchar(128);not null;check(role_id <> '')"`
	OrgId  string `json:"orgId" gorm:"column:org_id;type:varchar(128);not null;check(org_id <> '')"`
}
