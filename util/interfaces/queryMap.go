package interfaces

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type QueryMap map[string]interface{}

func NewQueryMap() QueryMap {
	return make(QueryMap)
}

func (i QueryMap) GetID() string {
	if id, ok := i["id"].(string); ok {
		return id
	}
	return ""
}

func (i QueryMap) GetIDs() []string {
	if ids, ok := i["ids"].([]string); ok {
		return ids
	}
	return nil
}

func (i QueryMap) GetCurrent() int {
	if current, ok := i["current"].(int); ok {
		return current
	}
	return 0
}

func (i QueryMap) GetPageSize() int {
	if pageSize, ok := i["pageSize"].(int); ok {
		return pageSize
	}
	return 10
}

func (i QueryMap) GetOffset() int {
	if offset, ok := i["offset"].(int); ok {
		return offset
	}
	return 0
}
func (i QueryMap) GetStarttime() int {
	if startTime, ok := i["startTime"].(int); ok {
		return startTime
	}
	return 0
}
func (i QueryMap) GetEndtime() int {
	if endtime, ok := i["endtime"].(int); ok {
		return endtime
	}
	return 0
}

func (i QueryMap) GetOrderBy() string {
	if orderBy, ok := i["orderBy"].(string); ok {
		return orderBy
	}
	return ""
}

func (i QueryMap) GetSortFlag() string {
	if sortFlag, ok := i["sortFlag"].(string); ok {
		return sortFlag
	}
	return ""
}

// 0 ：true & false  1: true   2: false
func (i QueryMap) GetDisabled() int {
	if disabled, ok := i["disabled"].(int); ok {
		return disabled
	}
	return 0
}

func (i QueryMap) Paginate(tableName string) func(db *gorm.DB) *gorm.DB {
	if tableName != "" {
		tableName += "."
	}
	return func(db *gorm.DB) *gorm.DB {
		page := i.GetCurrent()
		if page == 0 {
			page = 1
		}

		pageSize := i.GetPageSize()
		switch {
		case pageSize > 100:
			pageSize = 100
		}

		offset := (page - 1) * pageSize
		db.Offset(offset).Limit(pageSize)
		if strings.TrimSpace(i.GetOrderBy()) != "" {
			db.Order(strings.ReplaceAll(i.GetOrderBy(), ";", " "))
		} else {
			db.Order(fmt.Sprintf("%screatetime desc", tableName))
		}

		return db
	}
}

func (i QueryMap) BaseWhere(tableName ...string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var curTable string
		if len(tableName) > 0 {
			curTable = tableName[0] + "."
		}
		if i.GetIDs() != nil {
			db.Where(fmt.Sprintf("%sid=any(?)", curTable), i.GetIDs())
		}

		if i.GetStarttime() > 0 {
			db.Where(fmt.Sprintf("%screatetime>=?", curTable), i.GetStarttime())
		}
		if i.GetEndtime() > 0 {
			db.Where(fmt.Sprintf("%screatetime<=?", curTable), i.GetEndtime())
		}
		if i.GetDisabled() == 1 {
			db.Where(fmt.Sprintf("%sdisabled=true", curTable))
		} else if i.GetDisabled() == 2 {
			db.Where(fmt.Sprintf("%sdisabled=false", curTable))
		}
		return db
	}
}

func (i QueryMap) GetByKey(key string) interface{} {
	if v, ok := i[key]; ok {
		return v
	}
	return ""
}

func (i QueryMap) GetByKeyWithDefault(key string, df interface{}) interface{} {
	if v, ok := i[key]; ok {
		return v
	}
	return df
}
