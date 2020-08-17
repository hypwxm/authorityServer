package appController

import (
	"encoding/json"
	"github.com/hypwxm/rider"
	"babygrowing/config"
	"babygrowing/service/newsDynamicsComment/model"
	"babygrowing/service/newsDynamicsComment/service"
	"babygrowing/util/response"
)

func create(c rider.Context) {
	sender := response.NewSender()
	(func() {
		query := new(model.WbNewsDynamicsComment)
		err := json.Unmarshal(c.Body(), &query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		query.Publisher = c.GetLocals(config.AppUserTokenKey).(string)
		id, err := service.Create(query)
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
		entity.Publisher = c.GetLocals(config.AppUserTokenKey).(string)
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
		query.Publisher = c.GetLocals(config.AppUserTokenKey).(string)
		err = service.Del(query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success("")
	})()
	c.SendJson(200, sender)
}
