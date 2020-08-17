package controller

import (
	"encoding/json"
	"github.com/hypwxm/rider"
	"babygrowing/service/admin/rolePermission/source/model"
	"babygrowing/service/admin/rolePermission/source/service"
	"babygrowing/util/response"
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
		list, err := service.List(query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success(list)
	})()
	c.SendJson(200, sender)
}
