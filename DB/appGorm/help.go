package appGorm

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type BaseQuery struct {
	IDs       []string `json:"ids" db:"ids"`
	Current   int      `db:"current"`
	PageSize  int      `db:"pagesize"`
	Offset    int      `db:"offset"`
	Starttime int64    `db:"starttime"`
	Endtime   int64    `db:"endtime"`
	OrderBy   string   `json:"-" db:"order_by"`
	SortFlag  string   `db:"sort_flag"`
	Disabled  int      `db:"disabled"` // 0 ：true & false  1: true   2: false
}

func Paginate(query BaseQuery) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page := query.Current
		if page == 0 {
			page = 1
		}

		pageSize := query.PageSize
		switch {
		case pageSize > 100:
			pageSize = 100
		}

		offset := (page - 1) * pageSize
		db.Offset(offset).Limit(pageSize)
		if strings.TrimSpace(query.OrderBy) != "" {
			db.Order(strings.ReplaceAll(query.OrderBy, ";", " "))
		} else {
			db.Order("createtime desc")
		}

		return db
	}
}

func BaseWhere(query BaseQuery, tableName ...string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var curTable string
		if len(tableName) > 0 {
			curTable = tableName[0] + "."
		}
		if query.IDs != nil {
			db.Where(fmt.Sprintf("%sid=any(?)", curTable), query.IDs)
		}

		if query.Starttime > 0 {
			db.Where(fmt.Sprintf("%screatetime>=?", curTable), query.Starttime)
		}
		if query.Endtime > 0 {
			db.Where(fmt.Sprintf("%screatetime<=?", curTable), query.Endtime)
		}
		if query.Disabled == 1 {
			db.Where(fmt.Sprintf("%sdisabled=true", curTable))
		} else if query.Disabled == 2 {
			db.Where(fmt.Sprintf("%sdisabled=false", curTable))
		}
		return db
	}
}
