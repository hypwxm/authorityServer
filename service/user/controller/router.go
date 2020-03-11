package controller

import "github.com/hypwxm/rider"

func Router() *rider.Router {
	route := rider.NewRouter()

	// 获取注册的用户
	route.POST("/list", list)
	return route
}
