package pgsql

import (
	"fmt"
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

func BaseWhere(query BaseQuery, tableName ...string) string {
	/*if query == nil {
		query = new(BaseQuery)
	}*/

	var curTable string
	if len(tableName) > 0 {
		curTable = tableName[0] + "."
	}

	var whereSql = ""
	if query.IDs != nil {
		whereSql = whereSql + fmt.Sprintf(` and %sid=any(:ids)`, curTable)
	}

	if query.Starttime > 0 {
		whereSql = whereSql + fmt.Sprintf(` and %screatetime>=:starttime`, curTable)
	}
	if query.Endtime > 0 {
		whereSql = whereSql + fmt.Sprintf(` and %screatetime<=:endtime`, curTable)
	}
	if query.Disabled == 1 {
		whereSql = whereSql + fmt.Sprintf(` and %sdisabled=true`, curTable)
	} else if query.Disabled == 2 {
		whereSql = whereSql + fmt.Sprintf(` and %sdisabled=false`, curTable)
	}

	whereSql = whereSql + fmt.Sprintf(` and %sisdelete='false'`, curTable)
	return whereSql
}

func BaseOption(query BaseQuery, tableName ...string) string {
	/*if query == nil {
		return ""
	}*/
	var curTable string
	if len(tableName) > 0 {
		curTable = tableName[0] + "."
	}

	var optionSql string = ""
	if strings.TrimSpace(query.OrderBy) != "" {
		optionSql = optionSql + ` order by ` + query.OrderBy
	} else {
		optionSql = optionSql + fmt.Sprintf(` order by %screatetime desc`, curTable)
	}
	/*if strings.TrimSpace(query.SortFlag) != "" {
		optionSql = optionSql + ` ` + query.SortFlag
	} else {
		optionSql = optionSql + ` desc`
	}*/
	if query.Current > 0 {
		optionSql = optionSql + fmt.Sprintf(` limit %d`, query.PageSize)
	}
	query.Offset = (query.Current - 1) * query.PageSize

	if query.Offset > 0 {
		optionSql = optionSql + fmt.Sprintf(` offset %d`, query.Offset)
	}
	return optionSql
}
