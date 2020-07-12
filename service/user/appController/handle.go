package appController

import (
	"database/sql"
	"encoding/json"
	"github.com/hypwxm/rider"
	"worldbar/config"
	"worldbar/service/user/model"
	"worldbar/service/user/model/houseModel"
	"worldbar/service/user/service"
	"worldbar/util/response"
)

func registry(c rider.Context) {
	sender := response.NewSender()
	(func() {
		user := new(model.WbUser)
		err := json.Unmarshal(c.Body(), &user)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		user.Type = sql.NullString{
			String: "1",
			Valid:  true,
		}
		_, err = service.Create(user)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success("注册成功")
	})()
	c.SendJson(200, sender)
}

func modify(c rider.Context) {
	sender := response.NewSender()
	(func() {
		user := new(model.UpdateByIDQuery)
		err := json.Unmarshal(c.Body(), &user)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		userId := c.GetLocals(config.AppUserTokenKey).(string)
		user.ID = userId
		err = service.Modify(user)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success("修改成功")
	})()
	c.SendJson(200, sender)
}

func house(c rider.Context) {
	sender := response.NewSender()
	(func() {
		query := new(houseModel.GetQuery)
		query.UserId = c.GetLocals(config.AppUserTokenKey).(string)
		house, err := service.GetUserHouse(query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success(house)
	})()
	c.SendJson(200, sender)
}