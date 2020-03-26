package controller

import "github.com/hypwxm/rider"

func Router() *rider.Router {
	route := rider.NewRouter()

	// 生成枚举
	route.POST("/list", list)
	route.POST("/delete", del)
	route.POST("/toggleDisabled", toggleDisabled)
	route.POST("/get", get)
	return route
}
