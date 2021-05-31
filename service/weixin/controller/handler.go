package controller

import (
	"babygrow/service/weixin/config"
	"babygrow/service/weixin/model"
	"babygrow/service/weixin/service"

	"babygrow/util/response"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hypwxm/rider"
)

func getMiniUserInfo(c rider.Context) {
	sender := response.NewSender()
	(func() {
		var authEntity = new(model.MiniProgramUserAuth)
		err := json.Unmarshal(c.Body(), authEntity)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		authEntity.AppId = config.AppId
		authEntity.AppSecret = config.AppSecret
		infoByte, err := service.GetMiniUserInfo(authEntity)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		fmt.Printf("%s", infoByte)

		userInfo := new(model.MiniProgramUserInfo)
		if err = json.Unmarshal(infoByte, userInfo); err != nil {
			sender.Fail(err.Error())
			return
		}
		// 在这里可以处理保存小程序用户信息之类的操作
		if userId, err := service.StoreUserInfo(userInfo, true); err != nil {
			sender.Fail(err.Error())
			return
		} else {
			// 获取openid信息成功了，直接登录
			c.Jwt().Set("userId", userId)

			log.Println(c.Jwt().Get("userId"), 11111)
		}

		sender.Success(c.Jwt().GetToken())
	})()
	c.SendJson(200, sender)
}

// 前端传wx.login获取的code来获取微信的session_key
func code2sessionKey(c rider.Context) {
	sender := response.NewSender()
	(func() {
		var authEntity = new(model.MiniProgramUserAuth)
		body := c.Body()
		err := json.Unmarshal(body, authEntity)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		authEntity.AppId = config.AppId
		authEntity.AppSecret = config.AppSecret
		sessionKeyMap, err := service.Code2SessionKey(authEntity)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		// 缓存sessionKey
		if err = service.StorageSessionKey(authEntity.Code, sessionKeyMap.SessionKey); err != nil {
			sender.Fail(err.Error())
			return
		}

		// 如果session_key获取失败或者缓存失败了，需要前端删除之前的code，因为一个code只能获取一次session_key，
		// 也就是说，如果要更新session_ke，需要换一个新的code，code的过期时间和session_key的过期时间不搭嘎

		// sessionKeyMap中有用户openId信息，如果只是需要openId，不需要unionId的情况不需要去调上面getMiniUserInfo那步

		// 如果需要自定义存储信息，请修改这一步，这一步没有unionid,只有openid，需要unoinid可以在后续的授权中获取
		if err = StoreAndLogin(c, body, sessionKeyMap.Openid); err != nil {
			sender.Fail(err.Error())
			return
		}

		sender.Success(c.Jwt().GetToken())
	})()
	c.SendJson(200, sender)
}

func isSessionKeyValid(c rider.Context) {
	sender := response.NewSender()
	(func() {
		var authEntity = new(model.MiniProgramUserAuth)
		if err := json.Unmarshal(c.Body(), authEntity); err != nil {
			sender.Fail(err.Error())
			return
		}
		if _, err := service.GetCachedSessionKey(authEntity); err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success("缓存的sessionKey可用")
	})()
	c.SendJson(200, sender)
}

func StoreAndLogin(c rider.Context, userInfoBytes []byte, openId string) error {
	// 下面这一步去存储用户信息
	userInfo := new(model.MiniProgramUserInfo)
	if err := json.Unmarshal(userInfoBytes, &userInfo); err != nil {
		return err
	}
	userInfo.OpenId = openId
	if userId, err := service.StoreUserInfo(userInfo, false); err != nil {
		return err
	} else {
		// 获取openid信息成功了，直接登录
		c.Jwt().Set("userId", userId)
	}
	return nil
}
