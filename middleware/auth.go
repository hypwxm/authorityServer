package middleware

import (
	"babygrowing/config"
	"babygrowing/util/response"
	"log"
	"strings"

	"github.com/hypwxm/rider"
)

func Auth() rider.HandlerFunc {
	return func(c rider.Context) {
		user := c.Jwt().Get(config.AppServerTokenKey)
		token, _ := c.Jwt().Values()
		log.Println("admintokens", token)

		path := c.Path()
		// 开放性接口
		if strings.Contains(path, "/open/") {
			c.Next()
			return
		}

		sender := response.NewSender()
		// 登录的用户需要在jwt中存在role
		if userStr, ok := user.(string); ok {
			// val, err := redis.GetVal("user_" + userStr + "_" + c.Jwt().GetToken())
			if strings.TrimSpace(userStr) == "" {
				sender.Fail("未登录或登录已失效")
				sender.Code = 40301
				c.SendJson(200, sender)
				return
			}
			c.SetLocals(config.AppServerTokenKey, userStr)
			c.SetLocals(config.AppLoginUserName,  c.Jwt().Get(config.AppLoginUserName).(string))
			c.Next()

		} else {
			sender.Fail("请先登录")
			sender.Code = 40301
			c.SendJson(200, sender)
			return
		}
	}
}
