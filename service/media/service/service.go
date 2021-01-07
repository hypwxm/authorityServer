package service

import (
	"babygrowing/service/media/model"

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

func Del(query *model.DeleteQuery, tx ...*sqlx.Tx) error {
	return new(model.Media).Delete(query, tx...)
}

func InitMedias(list []*model.Media, businessName string, businessId string, creator string) []*model.Media {
	if len(list) == 0 {
		return nil
	}
	return model.InitMedias(list, businessName, businessId, creator)
}