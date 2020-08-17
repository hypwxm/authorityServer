package appController

import (
	"encoding/json"
	"github.com/hypwxm/rider"
	"babygrowing/service/matter/matter/model"
	"babygrowing/service/matter/matter/service"
	"babygrowing/util/response"
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
