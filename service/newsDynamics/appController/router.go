package appController

import "github.com/hypwxm/rider"

func Router() *rider.Router {
	route := rider.NewRouter()

	// 获取枚举列表
	route.POST("/open/list", list)
	// 修改枚举信息


	// 获取枚举列表
	route.POST("/open/enumsOptions", enumsOptions)

	route.POST("/open/getOptionDetail", getOptionDetail)


	// 获取关联信息
	route.POST("/open/getAssociates", getAssociates)

	// 根据用户的房屋信息获取详情
	route.POST("/getHouseDetail", getHouseDetail)


	return route
}
