package controller

import (
	"github.com/hypwxm/authorityServer/config"
	adminUserService "github.com/hypwxm/authorityServer/service/admin/user/service"
	"github.com/hypwxm/authorityServer/util/interfaces"
	"github.com/hypwxm/authorityServer/util/response"

	"github.com/hypwxm/rider"
)

type LoginForm struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

func login(c rider.Context) {
	sender := response.NewSender()
	(func() {
		var loginForm = interfaces.NewQueryMap()
		err := loginForm.FromByte(c.Body())
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		if loginForm.GetStringValue("password") == "" {
			if err != nil {
				sender.Fail("账号或密码错误")
				return
			}
		}
		user, err := adminUserService.Get(loginForm)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		c.Jwt().Set(config.AppServerTokenKey, user.GetID())
		c.Jwt().Set(config.AppLoginUserName, user.GetStringValue("username"))

		sender.Success(c.Jwt().GetToken())
	})()

	c.SendJson(200, sender)
}

func loginAdmin(c rider.Context) {
	sender := response.NewSender()
	(func() {
		var query = interfaces.NewQueryMap()
		err := query.FromByte(c.Body())
		query.Set("id", c.GetLocals(config.AppServerTokenKey).(string))
		user, err := adminUserService.Get(query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success(user)
	})()

	c.SendJson(200, sender)
}
