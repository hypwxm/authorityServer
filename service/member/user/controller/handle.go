package controller

import (
	"babygrow/config"
	familyMemberModel "babygrow/service/member/familyMember/model"
	familyMemberService "babygrow/service/member/familyMember/service"
	service "babygrow/service/member/user/service2"

	"babygrow/util/interfaces"
	"babygrow/util/response"
	"encoding/json"

	"github.com/hypwxm/rider"
)

func list(c rider.Context) {
	sender := response.NewSender()
	(func() {
		var query = interfaces.NewQueryMap()
		err := query.FromByte(c.Body())
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
		user := new(service.CreateModel)
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
		var query = interfaces.NewQueryMap()
		err := query.FromByte(c.Body())
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		query.Set("userId", c.GetLocals(config.MemberTokenKey))
		err = service.Modify(query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success("操作成功")
	})()
	c.SendJson(200, sender)
}

func modifyNickname(c rider.Context) {
	sender := response.NewSender()
	(func() {
		var query = interfaces.NewQueryMap()
		err := query.FromByte(c.Body())
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		query.Set("userId", c.GetLocals(config.MemberTokenKey))
		err = service.Modify(query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success("操作成功")
	})()
	c.SendJson(200, sender)
}

func modifyAvatar(c rider.Context) {
	sender := response.NewSender()
	(func() {
		var query = interfaces.NewQueryMap()
		err := query.FromByte(c.Body())
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		query.Set("userId", c.GetLocals(config.MemberTokenKey))
		err = service.ModifyAvatar(query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success("操作成功")
	})()
	c.SendJson(200, sender)
}

// 家庭邀请专用
func getInfoForFamilyInvite(c rider.Context) {
	sender := response.NewSender()
	(func() {
		var query = interfaces.NewQueryMap()
		err := query.FromByte(c.Body())
		if err != nil {
			sender.Fail(err.Error())
			return
		}

		// 要在家园中邀请成员，查询时必须带上成员的家园的id
		if query.GetStringValue("familyId") == "" {
			sender.Fail(err.Error())
			return
		}
		user, err := service.Get(query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		// 获取成员信息，判断该用户是否已经是该家园成员
		list, _, err := familyMemberService.List(&familyMemberModel.Query{
			FamilyId: query.GetStringValue("familyId"),
			UserId:   user.GetID(),
		})
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		var m = make(map[string]interface{})
		m["avatar"] = user.GetValue("avatar")
		m["id"] = user.GetID()
		m["account"] = user.GetValue("account")
		if len(list) == 0 {
			// 非成员，可以去邀请的
			m["isMember"] = false
		} else {
			// 已经是成员，不允许重复邀请
			m["isMember"] = true
		}
		sender.Success(m)
	})()
	c.SendJson(200, sender)
}
