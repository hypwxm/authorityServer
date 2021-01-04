package controller

import (
	"babygrowing/config"
	adminUserModel "babygrowing/service/admin/user/model"
	adminUserService "babygrowing/service/admin/user/service"
	"babygrowing/service/user/model"
	"babygrowing/service/user/service"
	"babygrowing/util"
	"babygrowing/util/database"
	"babygrowing/util/response"
	"encoding/json"

	"github.com/hypwxm/rider"
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
		if loginForm.Password == "" {
			if err != nil {
				sender.Fail("账号或密码错误")
				return
			}
		}
		user, err := adminUserService.GetUser(&adminUserModel.GAdminUser{
			Account:  loginForm.UserName,
			Password: loginForm.Password,
		})
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		c.Jwt().Set(config.AppServerTokenKey, user.ID)
		sender.Success(c.Jwt().GetToken())
	})()

	c.SendJson(200, sender)
}

func loginAdmin(c rider.Context) {
	sender := response.NewSender()
	(func() {
		query := new(adminUserModel.GAdminUser)
		query.ID = c.GetLocals(config.AppServerTokenKey).(string)
		user, err := adminUserService.GetUser(query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success(user)
	})()

	c.SendJson(200, sender)
}

func appLogin(c rider.Context) {
	sender := response.NewSender()
	(func() {
		user := new(model.WbUser)
		err := json.Unmarshal(c.Body(), &user)
		if err != nil {
			sender.Fail(util.ErrorFormat(err))
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
			sender.Fail(util.ErrorFormat(err))
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
			BaseColumns: database.BaseColumns{
				ID: c.GetLocals(config.AppUserTokenKey).(string),
			},
		})
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success(user)
	})()

	c.SendJson(200, sender)
}
