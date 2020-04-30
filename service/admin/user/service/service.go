package service

import (
	"worldbar/service/admin/user/model"
)

func Create(entity *model.WbAdminUser) (string, error) {
	return entity.Insert()
}

func Modify(updateQuery *model.UpdateByIDQuery) error {
	return new(model.WbAdminUser).Update(updateQuery)
}

func List(query *model.Query) ([]*model.ListModel, int64, error) {
	return new(model.WbAdminUser).List(query)
}

func Del(query *model.DeleteQuery) error {
	return new(model.WbAdminUser).Delete(query)
}

func Get(query *model.GetQuery) (*model.GetModel, error) {
	return new(model.WbAdminUser).GetByID(query)
}

func GetUser(query *model.WbAdminUser) (*model.WbAdminUser, error) {
	return new(model.WbAdminUser).Get(query)
}

// 默认创建一个超级管理员
func InitAdmin() {
	admin := &model.WbAdminUser{
		Account:  "admin",
		Password: "123456",
		Username: "管理员",
		Type:     1,
	}
	admin.Insert()
}
