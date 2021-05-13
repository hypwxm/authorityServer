package interfaces

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type QueryInterface interface {
	GetID() string
	GetIDs() pq.StringArray
	GetCurrent() int
	GetPageSize() int
	GetStarttime() int
	GetEndtime() int
	GetOrderBy() string
	GetSortFlag() string
	GetDisabled() int
	GetValue(key string) interface{}
	GetValueWithDefault(key string, df interface{}) interface{}
	Set(key string, value interface{})
	GetStringValue(key string) string
	ToStringArray(key string) pq.StringArray
	FromByte([]byte) error
}

type QueryMap map[string]interface{}

func NewQueryMap() QueryInterface {
	return make(QueryMap)
}

func (i QueryMap) GetID() string {
	if id, ok := i["id"].(string); ok {
		return id
	}
	return ""
}

func (i QueryMap) FromByte(b []byte) error {
	return json.Unmarshal(b, &i)
}

func (i QueryMap) GetIDs() pq.StringArray {
	if ids, ok := i["ids"].([]interface{}); ok {
		a := make([]string, len(ids))
		for k, v := range ids {
			if id, ok := v.(string); ok {
				a[k] = id
			}
		}
		return a
	}
	return nil
}

// 把interface的数组转成pq.stringArray，供sql用
func (i QueryMap) ToStringArray(key string) pq.StringArray {
	if ids, ok := i[key].(pq.StringArray); ok {
		return ids
	}
	if ids, ok := i[key].([]interface{}); ok {
		a := make([]string, len(ids))
		for k, v := range ids {
			if id, ok := v.(string); ok {
				a[k] = id
			}
		}
		return a
	}
	return nil
}

func (i QueryMap) GetCurrent() int {
	if current, ok := i["current"].(float64); ok {
		return int(current)
	}
	return 0
}

func (i QueryMap) GetPageSize() int {
	if pageSize, ok := i["pageSize"].(float64); ok {
		return int(pageSize)
	}
	return 10
}

func (i QueryMap) GetOffset() int {
	if offset, ok := i["offset"].(float64); ok {
		return int(offset)
	}
	return 0
}
func (i QueryMap) GetStarttime() int {
	if startTime, ok := i["startTime"].(float64); ok {
		return int(startTime)
	}
	return 0
}
func (i QueryMap) GetEndtime() int {
	if endtime, ok := i["endtime"].(float64); ok {
		return int(endtime)
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
	if disabled, ok := i["disabled"].(float64); ok {
		return int(disabled)
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

func (i QueryMap) GetValue(key string) interface{} {
	if v, ok := i[key]; ok {
		return v
	}
	return ""
}

func (i QueryMap) GetStringValue(key string) string {
	if s, ok := i.GetValue(key).(string); ok {
		return s
	}
	return ""
}

func (i QueryMap) Set(key string, value interface{}) {
	i[key] = value
}

func (i QueryMap) GetValueWithDefault(key string, df interface{}) interface{} {
	if v, ok := i[key]; ok {
		return v
	}
	return df
}
