package service

import (
	"github.com/hypwxm/authorityServer/service/media/model"
	"github.com/hypwxm/authorityServer/util/interfaces"

	"github.com/jmoiron/sqlx"
)

func Create(entity *model.Media, tx ...*sqlx.Tx) (string, error) {
	return entity.Insert(tx...)
}

func MultiCreate(list []*model.Media, tx ...*sqlx.Tx) error {
	if len(list) == 0 {
		return nil
	}
	return model.StoreMedias(list, tx...)
}

func List(query *model.Query) ([]*model.Media, int, error) {
	return new(model.Media).List(query)
}

// list为不定的结构体，通过反射的方式将拿到的媒体信息写入到list中每项的keyName对应的值和媒体的businessId相等的那个放到mediaName指定的结构体字段中
func ListWithMedia(ids []string, businessName string, list interface{}, mediaName string, keyName string) error {
	return new(model.Media).ListWithMedia(&model.Query{
		BusinessIds: ids,
		Businesses:  []string{businessName},
	}, list, mediaName, keyName)
}

// 类似上面的方法，但是为了方便使用，为map形式单独写个方法
func ListMapWithMedia(ids []string, businessName string, list interfaces.ModelMapSlice, mediaName string, keyName string) error {
	return new(model.Media).ListMapWithMedia(&model.Query{
		BusinessIds: ids,
		Businesses:  []string{businessName},
	}, list, mediaName, keyName)
}

// 为了减少代码重复，比如头像，我只需要查询到的第一张图片
func ListMapWithMediaFirst(ids []string, businessName string, list interfaces.ModelMapSlice, keyName string, valueKey string) error {
	return new(model.Media).ListMapWithMediaFirst(&model.Query{
		BusinessIds: ids,
		Businesses:  []string{businessName},
	}, list, keyName, valueKey)
}

func Del(query *model.DeleteQuery, tx ...*sqlx.Tx) error {
	return new(model.Media).Delete(query, tx...)
}

func InitMedias(list []*model.Media, businessName string, businessId string, creator string) []*model.Media {
	if len(list) == 0 {
		return nil
	}
	return model.InitMedias(list, businessName, businessId, creator)
}
