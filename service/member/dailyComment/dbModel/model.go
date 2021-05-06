package dbModel

import (
	"babygrow/DB/appGorm"
)

const BusinessName = "g_member_baby_grow_comment"

type GDailyComment struct {
	appGorm.BaseColumns

	Content string `json:"content" db:"content" gorm:"column:content;type:text;default '';not null"`

	UserId    string `json:"userId" db:"user_id" gorm:"column:user_id;type:varchar(128);not null;check(user_id <> '')"`
	BabyId    string `json:"babyId" db:"baby_id" gorm:"column:baby_id;type:varchar(128);not null;check(baby_id <> '')"`
	DiaryId   string `json:"diaryId" db:"diary_id" gorm:"column:diary_id;type:varchar(128);not null;check(diary_id <> '')"`
	CommentId string `json:"commentId" db:"comment_id" gorm:"column:comment_id;type:varchar(128);not null;default ''"`

	Sort int `json:"sort" db:"sort" gorm:"column:sort;not null;default 0"`
}
