package pgsql

import (
	"github.com/lib/pq"
	"strings"
)

type BaseQuery struct {
	IDs       pq.StringArray `db:"ids"`
	Current   int            `db:"current"`
	PageSize  int            `db:"pagesize"`
	Offset    int            `db:"offset"`
	Starttime int64          `db:"starttime"`
	Endtime   int64          `db:"endtime"`
	OrderBy   string         `db:"order_by"`
	SortFlag  string         `db:"sort_flag"`
	Disabled  int            `db:"disabled"` // 0 ：true & false  1: true   2: false
}

func BaseWhere(query BaseQuery) string {
	/*if query == nil {
		query = new(BaseQuery)
	}*/
	var whereSql = ""
	if query.IDs != nil {
		whereSql = whereSql + ` and id=any(:ids)`
	}

	if query.Starttime > 0 {
		whereSql = whereSql + ` and createtime>=:starttime`
	}
	if query.Endtime > 0 {
		whereSql = whereSql + ` and createtime<=:endtime`
	}
	if query.Disabled == 1 {
		whereSql = whereSql + ` and disabled=true`
	} else if query.Disabled == 2 {
		whereSql = whereSql + ` and disabled=false`
	}

	whereSql = whereSql + ` and isdelete='false'`
	return whereSql
}

func BaseOption(query BaseQuery) string {
	/*if query == nil {
		return ""
	}*/
	var optionSql string = ""
	if strings.TrimSpace(query.OrderBy) != "" {
		optionSql = optionSql + ` order by ` + query.OrderBy + `,createtime desc `
	} else {
		optionSql = optionSql + ` order by createtime desc`
	}
	/*if strings.TrimSpace(query.SortFlag) != "" {
		optionSql = optionSql + ` ` + query.SortFlag
	} else {
		optionSql = optionSql + ` desc`
	}*/
	if query.Current > 0 {
		optionSql = optionSql + ` limit :pagesize`
	}
	query.Offset = (query.Current - 1) * query.PageSize

	if query.Offset > 0 {
		optionSql = optionSql + ` offset :offset`
	}
	return optionSql
}
