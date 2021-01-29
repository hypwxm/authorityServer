package controller

import (
	"babygrowing/config"
	"babygrowing/service/admin/rolePermission/menu/model"
	"babygrowing/service/admin/rolePermission/menu/service"
	"babygrowing/util/response"
	"encoding/json"

	"github.com/hypwxm/rider"
)

func create(c rider.Context) {
	sender := response.NewSender()
	(func() {
		entity := new(model.SaveQuery)
		err := json.Unmarshal(c.Body(), &entity)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		id, err := service.Create(entity)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success(id)
	})()
	c.SendJson(200, sender)
}

func list(c rider.Context) {
	sender := response.NewSender()
	(func() {
		query := new(model.Query)
		err := json.Unmarshal(c.Body(), &query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		// 如果前段没传角色id过来，就根据当前的登录信息的角色来
		userId := c.GetLocals(config.AppServerTokenKey).(string)
		query.UserId = userId
		list, err := service.List(query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success(list)
	})()
	c.SendJson(200, sender)
}

func del(c rider.Context) {
	sender := response.NewSender()
	(func() {
		query := new(model.DeleteQuery)
		err := json.Unmarshal(c.Body(), &query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		err = service.Del(query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success("")
	})()
	c.SendJson(200, sender)
}
