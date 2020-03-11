package appController

import (
	"encoding/json"
	"github.com/hypwxm/rider"
	"worldbar/service/house/model/locationEnumsM"
	"worldbar/service/house/model/locationOptionsM"
	"worldbar/service/house/service"
	"worldbar/util/response"
)

func list(c rider.Context) {
	sender := response.NewSender()
	(func() {
		query := new(locationEnumsM.Query)
		err := json.Unmarshal(c.Body(), &query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		query.Disabled = 2
		list, total, err := service.EnumsList(query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.SuccessList(list, int(total))
	})()
	c.SendJson(200, sender)
}

func enumsOptions(c rider.Context) {
	sender := response.NewSender()
	(func() {
		query := new(locationOptionsM.Query)
		err := json.Unmarshal(c.Body(), &query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		query.Disabled = 2
		list, total, err := service.EnumsOptionsList(query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.SuccessList(list, int(total))
	})()
	c.SendJson(200, sender)
}

func getOptionDetail(c rider.Context) {
	sender := response.NewSender()
	(func() {
		query := new(locationOptionsM.GetQuery)
		err := json.Unmarshal(c.Body(), &query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		option, err := service.GetEnumsOption(query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success(option)
	})()
	c.SendJson(200, sender)
}

func getAssociates(c rider.Context) {
	sender := response.NewSender()
	(func() {
		query := new(locationOptionsM.AssociateGetQuery)
		err := json.Unmarshal(c.Body(), &query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		list, err := service.GetAssociates(query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success(list)
	})()
	c.SendJson(200, sender)
}

func getHouseDetail(c rider.Context) {
	sender := response.NewSender()
	(func() {
		query := new(locationOptionsM.GetListByIdsQuery)
		err := json.Unmarshal(c.Body(), &query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		list, err := service.GetOptionsByIds(query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success(list)
	})()
	c.SendJson(200, sender)
}
