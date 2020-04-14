package service

import (
	"worldbar/service/matter/matterElementOption/model"
)

func Create(entity *model.WbMatterElementOption) (string, error) {
	return entity.Insert()
}

func MulCreate(entitys []model.WbMatterElementOption) (string, error) {
	return new(model.WbMatterElementOption).MulInsert(entitys)
}


func Modify(updateQuery *model.UpdateByIDQuery) error {
	return new(model.WbMatterElementOption).Update(updateQuery)
}

func List(query *model.Query) ([]*model.ListModel, int64, error) {
	return new(model.WbMatterElementOption).List(query)
}

func Del(query *model.DeleteQuery) error {
	return new(model.WbMatterElementOption).Delete(query)
}

func Get(query *model.GetQuery) (*model.GetModel, error) {
	return new(model.WbMatterElementOption).GetByID(query)
}

func UpdateSort(query *model.UpdateSortQuery) error {
	return new(model.WbMatterElementOption).UpdateSort(query)
}

func UpdateStatus(query *model.UpdateStatusQuery) error {
	return new(model.WbMatterElementOption).UpdateStatus(query)
}
