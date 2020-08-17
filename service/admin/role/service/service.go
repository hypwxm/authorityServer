package service

import (
	"babygrowing/service/admin/role/model"
)

func Create(entity *model.WbAdminRole) (string, error) {
	return entity.Insert()
}

func Modify(updateQuery *model.UpdateByIDQuery) error {
	return new(model.WbAdminRole).Update(updateQuery)
}

func List(query *model.Query) ([]*model.ListModel, int64, error) {
	return new(model.WbAdminRole).List(query)
}

func Del(query *model.DeleteQuery) error {
	return new(model.WbAdminRole).Delete(query)
}

func Get(query *model.GetQuery) (*model.GetModel, error) {
	return new(model.WbAdminRole).GetByID(query)
}

func UpdateSort(query *model.UpdateSortQuery) error {
	return new(model.WbAdminRole).UpdateSort(query)
}

func UpdateStatus(query *model.UpdateStatusQuery) error {
	return new(model.WbAdminRole).UpdateStatus(query)
}
