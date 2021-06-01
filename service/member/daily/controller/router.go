package controller

import "github.com/hypwxm/rider"

func Router() *rider.Router {
	route := rider.NewRouter()

	// swagger:operation POST /member/daily/create create
	// ---
	// tags:
	//   - 日常
	// summary: 创建日记
	// description: 创建日记
	// Consumes:
	//   - application/json
	// Produces:
	//   - application/json
	// parameters:
	// - name: Authorization
	//   in: header
	//   description: token
	//   type: string
	//   required: true
	// - name: Body
	//   in: body
	//   schema:
	//     required:
	//       - date
	//     properties:
	//       date:
	//         description: 日记记录的日期
	//         type: string
	//         x-go-name: Date
	//       height:
	//         description: 身高
	//         type: float64
	//         x-go-name: Height
	// responses:
	//   '200':
	//     description: 成功
	//     properties:
	//       code:
	//         type: integer
	//       message:
	//         type: string
	//       data:
	//         type: string
	//   400: badReq
	route.POST("/create", create)

	// swagger:operation POST /member/daily/list list
	// ---
	// tags:
	//   - 日常
	// summary: 列表
	// description: 列表
	// Consumes:
	//   - application/json
	// Produces:
	//   - application/json
	// parameters:
	// - name: Authorization
	//   in: header
	//   description: token
	//   type: string
	//   required: true
	// - name: Body
	//   in: body
	//   schema:
	//     required:
	//       - babyId
	//     properties:
	//       babyId:
	//         description: 宝贝id
	//         type: string
	//         x-go-name: BabyId
	// responses:
	//   200: repoResp
	//   400: badReq
	route.POST("/list", list)
	route.POST("/modify", modify)
	route.POST("/delete", del)
	route.POST("/get", get)
	return route
}
