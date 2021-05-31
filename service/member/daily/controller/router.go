package controller

import "github.com/hypwxm/rider"

func Router() *rider.Router {
	route := rider.NewRouter()

	// swagger:operation POST /member/daily/create daily
	// ---
	// summary: 创建日记
	// description: 创建日记
	// parameters:
	// - name: Authorization
	//   in: header
	//   description: token
	//   type: string
	//   required: true
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
	route.POST("/create", create)
	route.POST("/list", list)
	route.POST("/modify", modify)
	route.POST("/delete", del)
	route.POST("/get", get)
	return route
}
