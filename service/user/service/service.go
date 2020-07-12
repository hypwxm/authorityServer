package service

import (
	"worldbar/service/user/model"
	"worldbar/service/user/model/houseModel"
)

func Create(user *model.WbUser) (string, error) {
	return user.Insert()
}

func List(query *model.Query) ([]*model.WbUser, int64, error) {
	return new(model.WbUser).List(query)
}

func GetUser(query *model.WbUser) (*model.WbUser, error) {
	return new(model.WbUser).Get(query)
}

func Modify(query *model.UpdateByIDQuery) error {
	return new(model.WbUser).Update(query)
}

// 获取用户房屋信息
func GetUserHouse(query *houseModel.GetQuery) ([]*houseModel.ListModel, error) {
	return new(houseModel.WbUserHouse).GetByUserId(query)
}
