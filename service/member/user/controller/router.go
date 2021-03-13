package controller

import "github.com/hypwxm/rider"

func Router() *rider.Router {
	route := rider.NewRouter()

	// 获取注册的用户
	route.POST("/list", list)
	route.POST("/toggleDisabled", toggleDisabled)

	// 用户注册
	route.POST("/open/register", registry)

	// 修改个人信息
	route.POST("/modify", modify)
	route.POST("/modifyNickname", modifyNickname)
	route.POST("/modifyAvatar", modifyAvatar)

	// 获取某个用户信息（家庭邀请用）
	route.POST("/getInfoForFamilyInvite", getInfoForFamilyInvite)

	return route
}
