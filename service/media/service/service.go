package service

import (
	"babygrowing/service/media/model"
)

func Create(entity *model.Media) (string, error) {
	return entity.Insert()
}

func MultiCreate(list []*model.Media) error {
	if len(list) == 0 {
		return nil
	}
	return model.StoreMedias(list)
}

func List(query *model.Query) ([]*model.ListModel, int64, error) {
	return new(model.Media).List(query)
}

func Del(query *model.DeleteQuery) error {
	return new(model.Media).Delete(query)
}

func InitMedias(list []*model.Media, businessName string, creator string) []*model.Media {
	if len(list) == 0 {
		return nil
	}
	return model.InitMedias(list, businessName, creator)
}
