package service

import (
	"babygrow/service/member/mybabies/model"
	userModel "babygrow/service/member/user/model"
	"log"
)

func Create(entity *model.GMyBabies) (string, error) {
	return entity.Insert()
}

func Modify(updateQuery *model.UpdateByIDQuery) error {
	return new(model.GMyBabies).Update(updateQuery)
}

func List(query *model.Query) ([]*model.ListModel, int64, error) {
	return new(model.GMyBabies).List(query)
}

func Del(query *model.DeleteQuery) error {
	return new(model.GMyBabies).Delete(query)
}

func Get(query *model.GetQuery) (*model.GetModel, error) {
	return new(model.GMyBabies).GetByID(query)
}

func ToggleDisabled(query *model.DisabledQuery) error {
	return new(model.GMyBabies).ToggleDisabled(query)
}

func GetBabyRelations(query *model.MbQuery) ([]*model.MbListModel, error) {
	if query.BabyId == "" {
		return nil, nil
	}
	return new(model.GMemberBabyRelation).List(query)
}

// 前段传过来的是账号，所以这里需要用账号去查询用户id
func CreateBabyRelations(query *model.GMemberBabyRelation) (string, error) {
	// 用账号名去查询用户信息，拿到用户id
	user, err := new(userModel.GMember).Get(&userModel.GMember{Account: query.Account})
	if err != nil {
		return "", err
	}
	query.UserId = user.ID

	log.Printf("%+v", query)
	return query.Insert(nil)
}
