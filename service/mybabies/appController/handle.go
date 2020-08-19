package appController

import (
	"babygrowing/service/newsDynamics/model"
	"babygrowing/service/newsDynamics/service"
	"babygrowing/util/response"
	"encoding/json"
	"github.com/hypwxm/rider"
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
		query.Status = 1
		list, total, err := service.List(query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.SuccessList(list, int(total))
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