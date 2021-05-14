package controller

import (
	"babygrow/config"
	service "babygrow/service/member/user/service2"
	"babygrow/util"
	"babygrow/util/interfaces"
	"babygrow/util/response"

	"github.com/hypwxm/rider"
)

type LoginForm struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

func memberLogin(c rider.Context) {
	sender := response.NewSender()
	(func() {
		query := interfaces.NewQueryMap()
		err := query.FromByte(c.Body())
		if err != nil {
			sender.Fail(util.ErrorFormat(err))
			return
		}
		if query.GetStringValue("account") == "" || query.GetStringValue("password") == "" {
			sender.Fail("账号或密码错误")
			return
		}
		query.Set("selects", "id,nickname,account,password,salt")
		user, err := service.Get(query)
		if err != nil {
			sender.Fail(util.ErrorFormat(err))
			return
		}
		c.Jwt().Set(config.MemberTokenKey, user.GetID())
		c.Jwt().Set(config.MemberLoginUserName, user.GetStringValue("nickname"))
		sender.Success(c.Jwt().GetToken())
	})()

	c.SendJson(200, sender)
}

func loginAppUser(c rider.Context) {
	sender := response.NewSender()
	(func() {
		query := interfaces.NewQueryMap()
		query.Set("id", c.GetLocals(config.MemberTokenKey).(string))
		query.Set("selects", "id,nickname,realname,gender,birthday,account")
		user, err := service.Get(query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success(user)
	})()

	c.SendJson(200, sender)
}
