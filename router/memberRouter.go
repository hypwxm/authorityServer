package router

import (
	baseControl "babygrow/base_control"
	"babygrow/middleware"
	controller2 "babygrow/service/member/auth/controller"
	daily "babygrow/service/member/daily/controller"
	dailyComment "babygrow/service/member/dailyComment/controller"
	family "babygrow/service/member/family/controller"
	familyMember "babygrow/service/member/familyMember/controller"
	message "babygrow/service/message/controller"

	controller3 "babygrow/service/member/mybabies/controller"
	"babygrow/service/member/user/controller"

	"github.com/hypwxm/rider"
)

func MemberRouter() *rider.Router {
	route := rider.NewRouter()
	route.USE(rider.RiderJwt("Authorization", "ni2QWN2asd23aw9d9j29j2d9aj9d!23", 10000000, nil, false))
	// 基础模块，图片上传等
	baseControl.Init(route)

	route.Kid("/auth", middleware.MemberAuth(), controller2.Router())

	route.Kid("/", middleware.MemberAuth(), controller.Router())

	route.Kid("/mybabies", middleware.MemberAuth(), controller3.Router())
	route.Kid("/babyGrow", middleware.MemberAuth(), daily.Router())
	route.Kid("/babyGrowComment", middleware.MemberAuth(), dailyComment.Router())
	route.Kid("/family", middleware.MemberAuth(), family.Router())
	route.Kid("/familyMember", middleware.MemberAuth(), familyMember.Router())
	route.Kid("/message", middleware.MemberAuth(), message.Router())

	return route
}
