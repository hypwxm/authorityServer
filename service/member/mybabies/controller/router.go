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

	route.POST("/relations", relations)
	route.POST("/createRelations", createRelations)
	route.POST("/delRelations", deleteRelations)

	/**
	* 通过邀请链接提交的申请，
	* 申请加入宝宝大家庭
	* params:
	  ------ roleName: 对应的角色
	  ------ userId: 需要登录，（为注册的用户需要先注册）
	*/
	route.POST("/applyJoinFamily", applyJoinFamily)

	return route
}
