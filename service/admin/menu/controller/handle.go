package controller

import (
	"babygrowing/config"
	"babygrowing/service/admin/menu/model"
	"babygrowing/service/admin/menu/service"
	"babygrowing/util/response"
	"encoding/json"

	"github.com/hypwxm/rider"
)

func create(c rider.Context) {
	sender := response.NewSender()
	(func() {
		entity := new(model.GMenu)
		err := json.Unmarshal(c.Body(), &entity)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		// if strings.TrimSpace(entity.ParentId) == "" {
		// 	sender.Fail("数据错误")
		// 	return
		// }
		entity.UserId = c.GetLocals(config.AppServerTokenKey).(string)
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
		entity := new(model.UpdateByIDQuery)
		err := json.Unmarshal(c.Body(), &entity)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		entity.UserId = c.GetLocals(config.AppServerTokenKey).(string)
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
		query := new(model.Query)
		err := json.Unmarshal(c.Body(), &query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		query.UserId = c.GetLocals(config.AppServerTokenKey).(string)
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