package dbModel

import (
	"babygrow/DB/appGorm"
)

const BusinessName = "g_member_baby"

type GMyBabies struct {
	appGorm.BaseColumns

	// 操作的用户
	UserID string `json:"userId" db:"user_id" gorm:"column:user_id;type:varchar(128);not null;check(user_id<>'');index"`
	// 姓名
	Name string `json:"name" db:"name" gorm:"column:name;type:varchar(50);not null;default ''"`
	// 生日,(公历生日)
	Birthday string `json:"birthday" db:"birthday" gorm:"column:birthday;type:varchar(50);not null;default ''"`
	// 性别
	Gender string `json:"gender" db:"gender" gorm:"column:gender;type:varchar(10);not null;default ''"`
	// 照片
	Avatar string `json:"avatar" db:"avatar" gorm:"column:avatar;type:varchar(500);not null;default ''"`
	// 身份证号
	IdCard string `json:"idCard" db:"id_card" gorm:"column:id_card;type:varchar(50);not null;default ''"`
	// 兴趣
	Hobby string `json:"hobby" db:"hobby" gorm:"column:hobby;type:text;not null;default ''"`
	// 特长
	GoodAt string `json:"goodAt" db:"good_at" gorm:"column:good_at;type:text;not null;default ''"`
	// 喜欢的食物
	FavoriteFood string `json:"favoriteFood" db:"favorite_food" gorm:"column:favorite_food;type:text;not null;default ''"`
	// 喜欢的颜色
	FavoriteColor string `json:"favoriteColor" db:"favorite_color" gorm:"column:favorite_color;type:text;not null;default ''"`
	// 志向
	Ambition string `json:"ambition" db:"ambition" gorm:"column:ambition;type:text;not null;default ''"`

	Weight float64 `json:"weight" db:"weight" gorm:"column:weight;not null;default 0"`
	Height float64 `json:"height" db:"height" gorm:"column:height;not null;default 0"`

	// 该字段新建表时不创建，主要是传给创建关系用的
	RoleName string `json:"roleName" db:"role_name" gorm:"-"`
}

const BusinessNameMB = "g_member_baby_relation"

type GMemberBabyRelation struct {
	appGorm.BaseColumns
	RoleName string `json:"roleName" db:"role_name" gorm:"column:role_name;type:varchar(10);not null;check(role_name <> '')"`
	BabyId   string `json:"babyId" db:"baby_id" gorm:"column:baby_id;type:varchar(128);not null;check(baby_id <> '');uniqueIndex:user_baby_id"`
	UserId   string `json:"userId" db:"user_id" gorm:"column:user_id;type:varchar(128);not null;check(user_id <> '');uniqueIndex:user_baby_id"`

	Account string `json:"account" db:"-" gorm:"-"`
}

const BusinessNameMBApply = "g_member_baby_relation_apply"

type GMemberBabyRelationApply struct {
	appGorm.BaseColumns
	RoleName  string `json:"roleName" db:"role_name" gorm:"column:role_name;type:varchar(10);not null;check(role_name <> '')"`
	BabyId    string `json:"babyId" db:"baby_id" gorm:"column:baby_id;type:varchar(128);not null;check(baby_id <> '');index;"`
	UserId    string `json:"userId" db:"user_id" gorm:"column:user_id;type:varchar(128);not null;check(user_id <> '');index;"`
	InviterId string `json:"inviterId" gorm:"column:inviter_id;type:varchar(128);not null;check(inviter_id <> '');index;"`
	Status    int    `json:"status" gorm:"column:status;type:smallint;not null;default 1;comment:1-申请中，2-同意，3-拒绝"`
}
