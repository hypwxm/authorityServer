package controller

import "github.com/hypwxm/rider"

/*
微信授权中心
*/
func Router() *rider.Router {
	route := rider.NewRouter()

	// 小程序用户信息
	route.POST("/miniprogram/userinfo", getMiniUserInfo)

	// 获取 sessionkey 存储 sessionKey
	route.POST("/miniprogram/code2sessionkey", code2sessionKey)

	// 前端需要判断微信给的session_key是否过期，同时也需要判断后台的session_key是否存在
	route.POST("/miniprogram/issessionkeyvalid", isSessionKeyValid)
	return route
}
