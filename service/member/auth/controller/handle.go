package controller

import (
	"babygrow/config"
	"babygrow/service/member/user/model"
	"babygrow/service/member/user/service"
	"babygrow/util"
	"babygrow/util/database"
	"babygrow/util/response"
	"encoding/json"

	"github.com/hypwxm/rider"
)

type LoginForm struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

func memberLogin(c rider.Context) {
	sender := response.NewSender()
	(func() {
		user := new(model.GMember)
		err := json.Unmarshal(c.Body(), &user)
		if err != nil {
			sender.Fail(util.ErrorFormat(err))
			return
		}
		if user.Account == "" || user.Password == "" {
			sender.Fail("账号或密码错误")
			return
		}
		user, err = service.GetUser(&model.GMember{
			Account:  user.Account,
			Password: user.Password,
		})
		if err != nil {
			sender.Fail(util.ErrorFormat(err))
			return
		}
		c.Jwt().Set(config.MemberTokenKey, user.ID)
		c.Jwt().Set(config.MemberLoginUserName, user.Nickname)
		sender.Success(c.Jwt().GetToken())
	})()

	c.SendJson(200, sender)
}

func loginAppUser(c rider.Context) {
	sender := response.NewSender()
	(func() {
		user, err := service.GetUser(&model.GMember{
			BaseColumns: database.BaseColumns{
				ID: c.GetLocals(config.MemberTokenKey).(string),
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
