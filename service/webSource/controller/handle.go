package controller

import (
	"encoding/json"
	"github.com/hypwxm/rider"
	"worldbar/service/webSource/model"
	"worldbar/service/webSource/service"
	"worldbar/util/response"
)

func create(c rider.Context) {
	sender := response.NewSender()
	(func() {
		entity := new(model.WbSettingsSource)
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

func modify(c rider.Context) {
	sender := response.NewSender()
	(func() {
		entity := new(model.UpdateByIDQuery)
		err := json.Unmarshal(c.Body(), &entity)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
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

func toggleDisabled(c rider.Context) {
	sender := response.NewSender()
	(func() {
		query := new(model.DisabledQuery)
		err := json.Unmarshal(c.Body(), &query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		err = service.ToggleDisabled(query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success("")
	})()
	c.SendJson(200, sender)
}