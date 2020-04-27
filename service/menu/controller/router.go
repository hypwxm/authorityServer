package controller

import "github.com/hypwxm/rider"

func Router() *rider.Router {
	route := rider.NewRouter()

	// 生成枚举
	route.POST("/create", create)
	route.POST("/list", list)
	route.POST("/modify", modify)
	route.POST("/delete", del)
	route.POST("/get", get)
	route.POST("/toggleDisabled", toggleDisabled)
	return route
}
