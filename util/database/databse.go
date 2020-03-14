package database

import "database/sql"

type BaseColumns struct {
	ID         string        `json:"id" db:"id"`
	Createtime int64         `json:"createtime" db:"createtime"`
	Updatetime sql.NullInt64 `json:"updatetime" db:"updatetime"`
	Deletetime sql.NullInt64 `json:"deletetime" db:"deletetime"`
	Isdelete   bool          `json:"isdelete" db:"isdelete"`
	Disabled   bool          `json:"disabled" db:"disabled"`
}