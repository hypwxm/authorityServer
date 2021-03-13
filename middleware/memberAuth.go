package middleware

import (
	"babygrow/config"
	"babygrow/util/response"
	"strings"

	"github.com/hypwxm/rider"
)

func MemberAuth() rider.HandlerFunc {
	return func(c rider.Context) {
		sender := response.NewSender()
		if c.Jwt() == nil {
			sender.Fail("未登录或登录已失效")
			sender.Code = 40301
			c.SendJson(200, sender)
			return
		}
		user := c.Jwt().Get(config.MemberTokenKey)

		path := c.Path()
		// 开放性接口
		if strings.Contains(path, "/open/") {
			c.Next()
			return
		}

		// 登录的用户需要在jwt中存在role
		if userStr, ok := user.(string); ok {
			// val, err := redis.GetVal("user_" + userStr + "_" + c.Jwt().GetToken())
			if strings.TrimSpace(userStr) == "" {
				sender.Fail("未登录或登录已过期")
				sender.Code = 40301
				c.SendJson(200, sender)
				return
			}
			c.SetLocals(config.MemberTokenKey, userStr)
			c.SetLocals(config.MemberLoginUserName, c.Jwt().Get(config.MemberLoginUserName).(string))

			c.SetLocals("userID", userStr)
			c.Next()

		} else {
			sender.Fail("请先登录")
			sender.Code = 40301
			c.SendJson(200, sender)
			return
		}
	}
}
