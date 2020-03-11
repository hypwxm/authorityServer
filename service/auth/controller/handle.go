package controller

import (
	"encoding/json"
	"github.com/hypwxm/rider"
	"worldbar/config"
	"worldbar/service/user/model"
	"worldbar/service/user/service"
	"worldbar/util/response"
)

type LoginForm struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

func login(c rider.Context) {
	sender := response.NewSender()
	(func() {
		var loginForm = new(LoginForm)
		err := json.Unmarshal(c.Body(), loginForm)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		if loginForm.UserName != "admin" || loginForm.Password != "123456" {
			sender.Fail("账号或密码错误")
			return
		}
		c.Jwt().Set(config.AppServerTokenKey, "admin")
		sender.Success(c.Jwt().GetToken())
	})()

	c.SendJson(200, sender)
}

func loginAdmin(c rider.Context) {
	sender := response.NewSender()
	(func() {
		sender.Success(map[string]interface{}{
			"id":   "21312",
			"name": "admin",
		})
	})()

	c.SendJson(200, sender)
}

func appLogin(c rider.Context) {
	sender := response.NewSender()
	(func() {
		user := new(model.WbUser)
		err := json.Unmarshal(c.Body(), &user)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		if user.Account == "" || user.Password == "" {
			sender.Fail("账号或密码错误")
			return
		}
		user, err = service.GetUser(&model.WbUser{
			Account:  user.Account,
			Password: user.Password,
		})
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		c.Jwt().Set(config.AppUserTokenKey, user.ID)
		sender.Success(c.Jwt().GetToken())
	})()

	c.SendJson(200, sender)
}

func loginAppUser(c rider.Context) {
	sender := response.NewSender()
	(func() {
		user, err := service.GetUser(&model.WbUser{
			ID: c.GetLocals(config.AppUserTokenKey).(string),
		})
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success(user)
	})()

	c.SendJson(200, sender)
}
