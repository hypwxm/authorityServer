package controller

import (
	"encoding/json"

	"github.com/hypwxm/authorityServer/config"
	"github.com/hypwxm/authorityServer/service/admin/user/model"
	"github.com/hypwxm/authorityServer/service/admin/user/service"
	"github.com/hypwxm/authorityServer/util/interfaces"
	"github.com/hypwxm/authorityServer/util/response"

	"github.com/hypwxm/rider"
)

func create(c rider.Context) {
	sender := response.NewSender()
	(func() {
		entity := new(service.CreateModel)
		err := json.Unmarshal(c.Body(), &entity)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		entity.Creator = c.GetLocals(config.AppLoginUserName).(string)
		entity.CreatorId = c.GetLocals(config.AppServerTokenKey).(string)
		id, err := service.Create(entity)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success(id)
	})()
	c.SendJson(200, sender)
}

func modify(c rider.Context) {
	sender := response.NewSender()
	(func() {
		entity := interfaces.NewQueryMap()
		err := entity.FromByte(c.Body())
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		entity.Set("userId", c.GetLocals(config.AppServerTokenKey).(string))
		err = service.Modify(entity)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success("操作成功")
	})()
	c.SendJson(200, sender)
}

func list(c rider.Context) {
	sender := response.NewSender()
	(func() {
		query := interfaces.NewQueryMap()
		err := query.FromByte(c.Body())
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		list, total, err := service.List(query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.SuccessList(list, int(total))
	})()
	c.SendJson(200, sender)
}

func del(c rider.Context) {
	sender := response.NewSender()
	(func() {
		query := interfaces.NewQueryMap()
		err := query.FromByte(c.Body())
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

func get(c rider.Context) {
	sender := response.NewSender()
	(func() {
		query := new(model.GetQuery)
		err := json.Unmarshal(c.Body(), &query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		entity, err := service.Get(query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success(entity)
	})()
	c.SendJson(200, sender)
}
