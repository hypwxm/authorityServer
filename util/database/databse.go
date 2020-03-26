package database

import (
	"database/sql"
	"worldbar/util"
)

type BaseColumns struct {
	ID         string        `json:"id" db:"id"`
	Createtime int64         `json:"createtime" db:"createtime"`
	Updatetime sql.NullInt64 `json:"updatetime" db:"updatetime"`
	Deletetime sql.NullInt64 `json:"deletetime" db:"deletetime"`
	Isdelete   bool          `json:"isdelete" db:"isdelete"`
	Disabled   bool          `json:"disabled" db:"disabled"`
}

func (s *BaseColumns) Init() {
	s.Createtime = util.GetCurrentMS()
	s.Isdelete = false
	s.Disabled = false
	s.ID = util.GetUuid()
}
