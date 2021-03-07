package service

import (
	"babygrow/service/member/user/model"
	"context"
)

func Create(user *model.GMember) (string, error) {
	return user.Insert()
}

func List(query *model.Query) ([]*model.GMember, int64, error) {
	return new(model.GMember).List(query)
}

// 客户端查询用，不会返回密码等隐私信息
func GetUserById(ctx context.Context, query *model.GetQuery) (*model.GetByIdModel, error) {
	return new(model.GMember).GetByID(ctx, query)
}

// 登录，验证密码时用，会返回密码等信息
func GetUser(query *model.GMember) (*model.GMember, error) {
	return new(model.GMember).Get(query)
}

func Modify(query *model.UpdateByIDQuery) error {
	return new(model.GMember).Update(query)
}

func ModifyNickname(query *model.UpdateByIDQuery) error {
	return new(model.GMember).UpdateNickname(query)
}

func ModifyAvatar(query *model.UpdateByIDQuery) error {
	return new(model.GMember).UpdateAvatar(query)
}

func ToggleDisabled(query *model.DisabledQuery) error {
	return new(model.GMember).ToggleDisabled(query)
}
