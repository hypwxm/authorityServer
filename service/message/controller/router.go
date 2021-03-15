package controller

import "github.com/hypwxm/rider"

func Router() *rider.Router {
	route := rider.NewRouter()

	// 生成枚举
	route.POST("/create", create)
	route.POST("/list", list)
	route.POST("/unreadcount", unreadcount)

	return route
}
