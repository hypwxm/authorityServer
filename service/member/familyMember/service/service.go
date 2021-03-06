package service

import (
	"babygrow/service/member/familyMember/model"
	"context"
)

func Create(ctx context.Context, entity *model.GFamilyMembers) (string, error) {
	return entity.Insert(ctx)
}

func Modify(updateQuery *model.UpdateByIDQuery) error {
	return new(model.GFamilyMembers).Update(updateQuery)
}

func List(query *model.Query) ([]*model.ListModel, int64, error) {
	return new(model.GFamilyMembers).List(query)
}

func Del(query *model.DeleteQuery) error {
	return new(model.GFamilyMembers).Delete(query)
}

func Get(query *model.GetQuery) (*model.GetModel, error) {
	return new(model.GFamilyMembers).GetByID(query)
}

func ToggleDisabled(query *model.DisabledQuery) error {
	return new(model.GFamilyMembers).ToggleDisabled(query)
}
