package appGorm

import (
	"authorityServer/util"

	"gorm.io/gorm"
)

type BaseColumns struct {
	ID         string         `json:"id" db:"id" gorm:"column:id;primaryKey;size:128"`
	Createtime int64          `json:"createtime" db:"createtime" gorm:"autoCreateTime:milli;column:createtime;not null"`
	Updatetime int64          `json:"updatetime" db:"updatetime" gorm:"autoUpdateTime:milli;column:updatetime;not null"`
	Deletetime int64          `json:"-" db:"deletetime" gorm:"-"`
	DeletetAt  gorm.DeletedAt `json:"-" db:"delete_at" gorm:"column:delete_at;index;"`
	Isdelete   bool           `json:"isdelete" db:"isdelete" gorm:"-"`
	Disabled   bool           `json:"disabled" db:"disabled" gorm:"column:disabled;default false;not null"`
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
