package appController

import "github.com/hypwxm/rider"

func Router() *rider.Router {
	route := rider.NewRouter()

	// 新评论
	route.POST("/create", create)

	// 分页获取评论列表
	route.POST("/list", list)

	// 编辑我的评论
	route.POST("/modify", modify)

	// 删除我的评论
	route.POST("/del", del)

	return route
}
