package service2

import (
	"babygrow/DB/appGorm"
	"babygrow/service/media/dao"
	"babygrow/service/media/dbModel"
	"babygrow/util/interfaces"
	"errors"
	"fmt"
	"strings"
)

type CreateModel struct {
	dbModel.Media
}

func Create(entity *CreateModel) (string, error) {
	if strings.TrimSpace(entity.UserID) == "" {
		return "", fmt.Errorf("操作错误")
	}
	if strings.TrimSpace(entity.Business) == "" {
		return "", fmt.Errorf("未知业务")
	}
	if strings.TrimSpace(entity.BusinessId) == "" {
		return "", fmt.Errorf("未知业务")
	}

	db := appGorm.Open()
	return dao.Insert(db, &entity.Media)
}

func Modify(query interfaces.QueryInterface) error {
	if query == nil {
		return errors.New("无更新条件")
	}
	if strings.TrimSpace(query.GetID()) == "" {
		return errors.New("更新条件错误")
	}
	db := appGorm.Open()
	return dao.Update(db, query)
}

func List(query interfaces.QueryInterface) (interfaces.ModelMapSlice, int64, error) {
	db := appGorm.Open()
	list, total, err := dao.List(db, query)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func Del(query interfaces.QueryInterface) error {
	db := appGorm.Open()
	return dao.Delete(db, query)
}

func Get(query interfaces.QueryInterface) (interfaces.ModelInterface, error) {
	db := appGorm.Open()
	return dao.Get(db, query)
}

// 业务查询媒体，把媒体对应的的每一项放到列表对应的每项中
func MergeMediaToListItem(query interfaces.QueryInterface, olist interfaces.ModelMapSlice, mediaName string, keyName string) error {
	medias, _, err := List(query)
	if err != nil {
		return err
	}
	if mediaName == "" {
		mediaName = "medias"
	}
	if keyName == "" {
		keyName = "id"
	}
	for _, v := range olist {
		v.Set(mediaName, make(interfaces.ModelMapSlice, 0))
		for _, vm := range medias {
			if v.GetValue(keyName) == vm.GetValue("businessId") {
				v.Set(mediaName, append(v.GetValue(mediaName).(interfaces.ModelMapSlice), vm))
			}
		}
	}
	return nil
}

// 业务查询媒体，把媒体对应的的第一项的媒体地址放到列表对应的每项中
func MergeFirstMediaToListItem(query interfaces.QueryInterface, olist interfaces.ModelMapSlice, keyName string, valueKey string) error {
	medias, _, err := List(query)
	if err != nil {
		return err
	}
	for _, v := range olist {
		for _, vm := range medias {
			if v.GetValue(keyName) == vm.GetValue("businessId") {
				v.Set(valueKey, vm.GetValue("url"))
				break
			}
		}
	}
	return nil
}

func InitMedias(list []*dbModel.Media, businessName string, businessId string, creator string) []*dbModel.Media {
	if len(list) == 0 {
		return nil
	}
	medias := make([]*dbModel.Media, 0)
	for _, v := range list {
		if v == nil {
			continue
		}
		// 有一种情况是这个媒体是之前的信息，就不进行重新保存了
		if !(v.Business != "" && v.BusinessId != "" && v.UserID != "" && v.ID != "") {
			v.Business = businessName
			v.BusinessId = businessId
			v.UserID = creator
			medias = append(medias, v)
		}
	}
	return medias
}

func MultiCreate(list []*dbModel.Media) error {
	if len(list) == 0 {
		return nil
	}
	db := appGorm.Open()
	for _, v := range list {
		if v == nil {
			return fmt.Errorf("文件错误")
		}
		_, err := dao.Insert(db, v)
		if err != nil {
			return err
		}
	}
	return nil
}
