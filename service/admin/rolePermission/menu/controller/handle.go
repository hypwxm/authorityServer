package controller

import (
	"github.com/hypwxm/authorityServer/config"
	"github.com/hypwxm/authorityServer/service/admin/rolePermission/menu/service"
	"github.com/hypwxm/authorityServer/util/interfaces"
	"github.com/hypwxm/authorityServer/util/response"

	"github.com/hypwxm/rider"
)

func create(c rider.Context) {
	sender := response.NewSender()
	(func() {
		entity := interfaces.NewQueryMap()
		err := entity.FromByte(c.Body())
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		id, err := service.Save(entity)
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
		query := interfaces.NewQueryMap()
		err := query.FromByte(c.Body())
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		// 如果前段没传角色id过来，就根据当前的登录信息的角色来
		query.Set("userId", c.GetLocals(config.AppServerTokenKey).(string))
		list, err := service.List(query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success(list)
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
		err = service.Del(query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success("")
	})()
	c.SendJson(200, sender)
}
