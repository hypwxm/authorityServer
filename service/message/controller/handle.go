package controller

import (
	"babygrow/service/message/model"
	"babygrow/service/message/service"
	"babygrow/util/response"
	"encoding/json"

	"github.com/hypwxm/rider"
)

func create(c rider.Context) {
	sender := response.NewSender()
	(func() {
		entity := new(model.GMessage)
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
		list, total, err := service.List(query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.SuccessList(list, int(total))
	})()
	c.SendJson(200, sender)
}

func unreadcount(c rider.Context) {
	sender := response.NewSender()
	(func() {
		query := new(model.UnreadCountQuery)
		err := json.Unmarshal(c.Body(), &query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		total, err := service.UnreadCount(query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success(int(total))
	})()
	c.SendJson(200, sender)
}
