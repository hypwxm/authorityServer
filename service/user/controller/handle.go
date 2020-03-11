package controller

import (
	"encoding/json"
	"github.com/hypwxm/rider"
	"worldbar/service/user/model"
	"worldbar/service/user/service"
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
		list, total, err := service.List(query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.SuccessList(list, int(total))
	})()
	c.SendJson(200, sender)
}
