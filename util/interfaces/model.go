package interfaces

import "strings"

type ModelInterface interface {
	GetID() string
}

type ModelMap map[string]interface{}

type ModelMapSlice []ModelMap

func NewModelMap() ModelMap {
	return make(ModelMap)
}

func NewModelMapFromMap(m map[string]interface{}) ModelMap {
	nm := make(ModelMap)
	for k, v := range m {
		nm[k] = v
	}
	return nm
}

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

func (i ModelMap) GetID() string {
	if id, ok := i["id"].(string); ok {
		return id
	}
	return ""
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
		m[s] = v
	}
	return m
}

func (is ModelMapSlice) ToCamelKey() ModelMapSlice {
	for k, v := range is {
		nv := v.ToCamelKey()
		is[k] = nv
	}
	return is
}
