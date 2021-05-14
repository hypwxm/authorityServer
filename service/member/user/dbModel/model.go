package dbModel

import "babygrow/util/database"

const BusinessName = "member"

type GMember struct {
	database.BaseColumns

	Nickname string `json:"nickname" db:"nickname" gorm:"column:nickname;type:varchar(40);not null;default ''"`
	RealName string `json:"realName" db:"realname" gorm:"column:realname;type:varchar(40);not null;default ''"`
	Gender   string `json:"gender" gorm:"column:gender;type:varchar(40);not null;default ''"`

	Birthday string `json:"birthday" gorm:"column:birthday;type:varchar(40);not null;default ''"`
	Phone    string `json:"phone" db:"phone" gorm:"column:phone;type:varchar(40);not null;default ''"`
	Account  string `json:"account" db:"account" gorm:"column:account;type:varchar(40);not null;check(account <> '');uniqueIndex"`
	Password string `json:"password" db:"password"`
	Salt     string `json:"-" db:"salt" gorm:"column:salt;type:varchar(40);not null;check(salt <> '')"`
}
