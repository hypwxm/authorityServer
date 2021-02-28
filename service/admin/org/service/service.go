package service

import (
	"babygrow/service/admin/org/model"
)

func Create(entity *model.GOrg) (string, error) {
	return entity.Insert()
}

func Modify(updateQuery *model.UpdateByIDQuery) error {
	return new(model.GOrg).Update(updateQuery)
}

func List(query *model.Query) ([]*model.ListModel, int64, error) {
	return new(model.GOrg).List(query)
}

func Del(query *model.DeleteQuery) error {
	return new(model.GOrg).Delete(query)
}

func Get(query *model.GetQuery) (*model.GetModel, error) {
	return new(model.GOrg).GetByID(query)
}
