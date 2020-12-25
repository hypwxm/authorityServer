package service

import (
	"babygrowing/service/admin/user/model"
)

func Create(entity *model.GAdminUser) (string, error) {
	return entity.Insert()
}

func Modify(updateQuery *model.UpdateByIDQuery) error {
	return new(model.GAdminUser).Update(updateQuery)
}

func List(query *model.Query) ([]*model.ListModel, int64, error) {
	return new(model.GAdminUser).List(query)
}

func Del(query *model.DeleteQuery) error {
	return new(model.GAdminUser).Delete(query)
}

func Get(query *model.GetQuery) (*model.GetModel, error) {
	return new(model.GAdminUser).GetByID(query)
}

func GetUser(query *model.GAdminUser) (*model.GAdminUser, error) {
	return new(model.GAdminUser).Get(query)
}

// 默认创建一个超级管理员
func InitAdmin() {
	admin := &model.GAdminUser{
		Account:  "admin",
		Password: "123456",
		Username: "管理员",
	}
	admin.Insert()
}
