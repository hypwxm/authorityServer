package model

import (
	"babygrow/DB/appGorm"
	mediaModel "babygrow/service/media/model"
)

type GDaily struct {
	appGorm.BaseColumns

	// 今日份体重
	Weight float64 `json:"weight" db:"weight" gorm:"column:weight;not null;default 0"`
	// 今日份身高
	Height float64 `json:"height" db:"height" gorm:"column:height;not null;default 0"`

	// 今日份记录
	Diary string `json:"diary" db:"diary" gorm:"column:diary;type:text;not null;default ''"`

	UserId string `json:"userId" db:"user_id" gorm:"column:user_id;type:varchar(128);not null;check(user_id <> '')"`
	BabyId string `json:"babyId" db:"baby_id" gorm:"column:baby_id;type:varchar(128);not null;check(baby_id <> '')"`

	Weather     string  `json:"weather" db:"weather" gorm:"column:weather;type:varchar(50);not null;default ''"`
	Mood        string  `json:"mood" db:"mood" gorm:"column:mood;type:varchar(40);not null;default ''"`
	Health      string  `json:"health" db:"health" gorm:"column:health;type:varchar(40);not null;default ''"`
	Temperature float64 `json:"temperature" db:"temperature" gorm:"column:temperature;not null;default 0"`

	Date string `json:"date" db:"date" gorm:"column:date;type:varchar(40);not null;default ''"`

	Sort int `json:"sort" db:"sort" gorm:"column:sort;not null;default 0"`

	Medias []*mediaModel.Media `json:"medias" gorm:"-"`
}
