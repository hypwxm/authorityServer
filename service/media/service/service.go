package service

import (
	"babygrow/service/media/model"

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

func ListWithMedia(ids []string, businessName string, list interface{}, mediaName string, keyName string) error {
	return new(model.Media).ListWithMedia(&model.Query{
		BusinessIds: ids,
		Businesses:  []string{businessName},
	}, list, mediaName, keyName)
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
