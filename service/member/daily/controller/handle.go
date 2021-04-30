package controller

import (
	"babygrow/config"
	"babygrow/service/member/daily/service"
	"babygrow/util/interfaces"
	"babygrow/util/response"
	"encoding/json"
	"log"

	"github.com/hypwxm/rider"
)

func create(c rider.Context) {
	sender := response.NewSender()
	(func() {
		entity := new(service.CreateModel)
		err := json.Unmarshal(c.Body(), &entity)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		entity.UserId = c.GetLocals(config.MemberTokenKey).(string)
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
		entity := interfaces.NewQueryMap()
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
		var query = interfaces.NewQueryMap()
		err := query.FromByte(c.Body())
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		query.Set("userId", c.GetLocals(config.MemberTokenKey).(string))
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
		query := interfaces.NewQueryMap()
		err := query.FromByte(c.Body())
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		log.Println(query.GetIDs(), len(query.GetIDs()))
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
		query := interfaces.NewQueryMap()
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
