package controller

import (
	"github.com/hypwxm/rider"
)

func Router() *rider.Router {
	route := rider.NewRouter()

	// swagger:operation POST /member/open/login repos login
	// ---
	// summary: app用户登录
	// description: app用户登录
	// parameters:
	// - name: account
	//   in: body
	//   description: 用户名
	//   type: string
	//   required: true
	// - name: password
	//   in: body
	//   description: 登录密码
	//   type: string
	//   required: true
	// responses:
	//   200: repoResp
	//   400: badReq
	route.POST("/open/login", memberLogin)
	route.POST("/loginUser", loginAppUser)

	return route
}
