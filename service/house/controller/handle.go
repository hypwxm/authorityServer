package controller

import (
	"encoding/json"
	"github.com/hypwxm/rider"
	"worldbar/service/house/model/locationEnumsM"
	"worldbar/service/house/model/locationOptionsM"
	"worldbar/service/house/service"
	"worldbar/util/response"
)

func createEnums(c rider.Context) {
	sender := response.NewSender()
	(func() {
		entity := new(locationEnumsM.WbHouseEnums)
		err := json.Unmarshal(c.Body(), &entity)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		id, err := service.CreateEnums(entity)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success(id)
	})()
	c.SendJson(200, sender)
}

func modifyEnums(c rider.Context) {
	sender := response.NewSender()
	(func() {
		entity := new(locationEnumsM.UpdateByIDQuery)
		err := json.Unmarshal(c.Body(), &entity)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		err = service.ModifyEnums(entity)
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
		query := new(locationEnumsM.Query)
		err := json.Unmarshal(c.Body(), &query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		list, total, err := service.EnumsList(query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.SuccessList(list, int(total))
	})()
	c.SendJson(200, sender)
}

func deleteEnums(c rider.Context) {
	sender := response.NewSender()
	(func() {
		query := new(locationEnumsM.DeleteQuery)
		err := json.Unmarshal(c.Body(), &query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		err = service.DeleteEnums(query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success("")
	})()
	c.SendJson(200, sender)
}

func updateSort(c rider.Context) {
	sender := response.NewSender()
	(func() {
		query := new(locationEnumsM.UpdateSortQuery)
		err := json.Unmarshal(c.Body(), &query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		err = service.UpdateSort(query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success("")
	})()
	c.SendJson(200, sender)
}




func createEnumsOption(c rider.Context) {
	sender := response.NewSender()
	(func() {
		entity := new(locationOptionsM.WbHouseOption)
		err := json.Unmarshal(c.Body(), &entity)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		id, err := service.CreateEnumsOption(entity)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success(id)
	})()
	c.SendJson(200, sender)
}

func modifyEnumsOption(c rider.Context) {
	sender := response.NewSender()
	(func() {
		entity := new(locationOptionsM.UpdateByIDQuery)
		err := json.Unmarshal(c.Body(), &entity)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		err = service.ModifyEnumsOption(entity)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success("操作成功")
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
		list, total, err := service.EnumsOptionsList(query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.SuccessList(list, int(total))
	})()
	c.SendJson(200, sender)
}

func deleteEnumsOption(c rider.Context) {
	sender := response.NewSender()
	(func() {
		query := new(locationOptionsM.DeleteQuery)
		err := json.Unmarshal(c.Body(), &query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		err = service.DeleteEnumsOption(query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success("")
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

func associate(c rider.Context) {
	sender := response.NewSender()
	(func() {
		query := new(locationOptionsM.AssociateQuery)
		err := json.Unmarshal(c.Body(), &query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		err = service.Associate(query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success("")
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

func deleteAssociates(c rider.Context) {
	sender := response.NewSender()
	(func() {
		query := new(locationOptionsM.AssociateQuery)
		err := json.Unmarshal(c.Body(), &query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		err = service.DeleteAssociates(query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success("")
	})()
	c.SendJson(200, sender)
}
