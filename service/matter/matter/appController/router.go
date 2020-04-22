package appController

import "github.com/hypwxm/rider"

func Router() *rider.Router {
	route := rider.NewRouter()

	// 获取处于发布状态的事宜
	route.POST("/list", list)

	return route
}
