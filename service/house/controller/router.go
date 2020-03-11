package controller

import "github.com/hypwxm/rider"

func Router() *rider.Router {
	route := rider.NewRouter()

	// 生成枚举
	route.POST("/createEnums", createEnums)
	// 获取枚举列表
	route.POST("/list", list)
	// 修改枚举信息
	route.POST("/modifyEnums", modifyEnums)

	route.POST("/deleteEnums", deleteEnums)
	route.POST("/updateSort", updateSort)


	// 生成枚举
	route.POST("/createEnumsOption", createEnumsOption)
	// 获取枚举列表
	route.POST("/enumsOptions", enumsOptions)
	// 修改枚举信息
	route.POST("/modifyEnumsOption", modifyEnumsOption)

	route.POST("/deleteEnumsOption", deleteEnumsOption)

	route.POST("/getOptionDetail", getOptionDetail)

	// 对属性进行关联
	route.POST("/associate", associate)

	// 获取关联信息
	route.POST("/getAssociates", getAssociates)

	route.POST("/deleteAssociates", deleteAssociates)

	return route
}
