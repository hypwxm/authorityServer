package interfaces

import (
	"log"
	"strings"
)

type ModelInterface interface {
	GetID() string
	GetValue(key string) interface{}
	Set(key string, value interface{})
	GetStringValue(key string) string
	GetValueWithDefault(key string, df interface{}) interface{}
	// ToCamelKey() ModelInterface
}

type ModelMap map[string]interface{}

func NewModelMap() ModelMap {
	return make(ModelMap)
}

func init() {
	var a ModelInterface = NewModelMap()
	log.Print(a)
}

func NewModelMapFromMap(m map[string]interface{}) ModelMap {
	nm := make(ModelMap)
	for k, v := range m {
		nm[k] = v
	}
	return nm
}

func (i ModelMap) GetID() string {
	if id, ok := i["id"].(string); ok {
		return id
	}
	return ""
}

func (i ModelMap) GetValue(key string) interface{} {
	if v, ok := i[key]; ok {
		return v
	}
	return ""
}

func (i ModelMap) GetStringValue(key string) string {
	if s, ok := i.GetValue(key).(string); ok {
		return s
	}
	return ""
}

func (i ModelMap) Set(key string, value interface{}) {
	i[key] = value
}

func (i ModelMap) GetValueWithDefault(key string, df interface{}) interface{} {
	if v, ok := i[key]; ok {
		return v
	}
	return df
}

func (i ModelMap) ToCamelKey() ModelMap {
	var m = NewModelMap()
	for k, v := range i {
		temp := strings.Split(k, "_")
		var s string
		for k, sv := range temp {
			vv := []rune(sv)
			if k == 0 {
				s += string(vv)
				continue
			}
			if len(vv) > 0 {
				if bool(vv[0] >= 'a' && vv[0] <= 'z') { //首字母大写
					vv[0] -= 32
				}
				s += string(vv)
			}
		}
		m.Set(s, v)
	}
	return m
}

type ModelMapSlice []ModelMap

func NewModelMapSlice(cap int) ModelMapSlice {
	return make(ModelMapSlice, cap)
}

func NewModelMapSliceFromMapSlice(s []map[string]interface{}) ModelMapSlice {
	nm := make([]ModelMap, 0)
	for k := range s {
		nm = append(nm, NewModelMapFromMap(s[k]))
	}
	return nm
}

func (is ModelMapSlice) ToCamelKey() ModelMapSlice {
	for k, v := range is {
		nv := v.ToCamelKey()
		is[k] = nv
	}
	return is
}

// 从slice的每项中取出key对应的value，业务中可以知道对应的数据类型，通过断言拿到具体类型
func (is ModelMapSlice) GetValues(key string) interface{} {
	var list = make([]interface{}, len(is))
	for _, v := range is {
		list = append(list, v[key])
	}
	return list
}

func (is ModelMapSlice) GetStringValues(key string) []string {
	var list = make([]string, len(is))
	for _, v := range is {
		list = append(list, v[key].(string))
	}
	return list
}
