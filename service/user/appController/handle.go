package appController

import (
	"babygrowing/config"
	"babygrowing/service/user/model"
	"babygrowing/service/user/service"
	"babygrowing/util/response"
	"encoding/json"

	"github.com/hypwxm/rider"
)

func registry(c rider.Context) {
	sender := response.NewSender()
	(func() {
		user := new(model.GMember)
		err := json.Unmarshal(c.Body(), &user)
		if err != nil {
			sender.Fail(err.Error())
			return
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
