package controller

import (
	"babygrow/config"
	"babygrow/service/member/user/model"
	"babygrow/service/member/user/service"
	"babygrow/util/response"
	"encoding/json"

	"github.com/hypwxm/rider"
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
		list, total, err := service.List(query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.SuccessList(list, int(total))
	})()
	c.SendJson(200, sender)
}

func registry(c rider.Context) {
	sender := response.NewSender()
	(func() {
		user := new(model.GMember)
		err := json.Unmarshal(c.Body(), &user)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		_, err = service.Create(user)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success("注册成功")
	})()
	c.SendJson(200, sender)
}

func modify(c rider.Context) {
	sender := response.NewSender()
	(func() {
		user := new(model.UpdateByIDQuery)
		err := json.Unmarshal(c.Body(), &user)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		userId := c.GetLocals(config.MemberTokenKey).(string)
		user.ID = userId
		err = service.Modify(user)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success("修改成功")
	})()
	c.SendJson(200, sender)
}

func modifyNickname(c rider.Context) {
	sender := response.NewSender()
	(func() {
		user := new(model.UpdateByIDQuery)
		err := json.Unmarshal(c.Body(), &user)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		userId := c.GetLocals(config.MemberTokenKey).(string)
		user.ID = userId
		err = service.ModifyNickname(user)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success("修改成功")
	})()
	c.SendJson(200, sender)
}

func modifyAvatar(c rider.Context) {
	sender := response.NewSender()
	(func() {
		user := new(model.UpdateByIDQuery)
		err := json.Unmarshal(c.Body(), &user)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		userId := c.GetLocals(config.MemberTokenKey).(string)
		user.ID = userId
		err = service.ModifyAvatar(user)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success("修改成功")
	})()
	c.SendJson(200, sender)
}

func toggleDisabled(c rider.Context) {
	sender := response.NewSender()
	(func() {
		query := new(model.DisabledQuery)
		err := json.Unmarshal(c.Body(), &query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		err = service.ToggleDisabled(query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success("")
	})()
	c.SendJson(200, sender)
}

// 家庭邀请专用
func getInfoForFamilyInvite(c rider.Context) {
	sender := response.NewSender()
	(func() {
		query := new(model.GMember)
		err := json.Unmarshal(c.Body(), &query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		user, err := service.GetUser(query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		var m = make(map[string]interface{})
		m["avatar"] = user.Avatar
		m["id"] = user.ID
		m["account"] = user.Account
		sender.Success(m)
	})()
	c.SendJson(200, sender)
}
