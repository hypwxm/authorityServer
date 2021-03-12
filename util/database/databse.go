package database

import (
	"babygrow/util"

	"gorm.io/gorm"
)

type BaseColumns struct {
	ID         string         `json:"id" db:"id" gorm:"column:id;primaryKey;size:128"`
	Createtime int64          `json:"createtime" db:"createtime" gorm:"autoUpdateTime:milli;column:createtime"`
	Updatetime int64          `json:"up" db:"updatetime" gorm:"autoUpdateTime:milli;column:updatetime"`
	Deletetime gorm.DeletedAt `json:"-" db:"deletetime" gorm:"column:deletetime;index"`
	Isdelete   bool           `json:"isdelete" db:"isdelete" gorm:"column:isdelete"`
	Disabled   bool           `json:"disabled" db:"disabled" gorm:"column:disabled"`
}

func (s *BaseColumns) Init() {
	s.Createtime = util.GetCurrentMS()
	s.Isdelete = false
	s.ID = util.GetUuid()
}

type BaseIDColumns struct {
	ID string `json:"id" db:"id"`
}

func (s *BaseIDColumns) Init() {
	s.ID = util.GetUuid()
}
