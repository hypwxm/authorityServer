package database

import (
	"babygrowing/util"
	"database/sql"
)

type BaseColumns struct {
	ID         string        `json:"id" db:"id"`
	Createtime int64         `json:"createtime" db:"createtime"`
	Updatetime sql.NullInt64 `json:"-" db:"updatetime"`
	Deletetime sql.NullInt64 `json:"-" db:"deletetime"`
	Isdelete   bool          `json:"isdelete" db:"isdelete"`
	Disabled   bool          `json:"disabled" db:"disabled"`
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
