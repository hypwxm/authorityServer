package appController

import (
	"encoding/json"
	"github.com/hypwxm/rider"
	"worldbar/config"
	"worldbar/service/like/model"
	"worldbar/service/like/service"
	"worldbar/util/response"
)

func list(c rider.Context) {
	sender := response.NewSender()
	(func() {
		query := new(model.Query)
		err := json.Unmarshal(c.Body(), &query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		query.Disabled = 2
		list, total, err := service.List(query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.SuccessList(list, int(total))
	})()
	c.SendJson(200, sender)
}

func create(c rider.Context) {
	sender := response.NewSender()
	(func() {
		query := new(model.WbLike)
		err := json.Unmarshal(c.Body(), &query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		query.UserId = c.GetLocals(config.AppUserTokenKey).(string)
		id, err := service.Create(query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success(id)
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
