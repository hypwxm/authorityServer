package router

import (
	"babygrowing/middleware"
	controller2 "babygrowing/service/member/auth/controller"
	daily "babygrowing/service/member/daily/controller"
	controller3 "babygrowing/service/member/mybabies/controller"
	"babygrowing/service/member/user/controller"

	"github.com/hypwxm/rider"
)

func MemberRouter() *rider.Router {
	route := rider.NewRouter()
	route.USE(rider.RiderJwt("Authorization", "ni2QWN2asd23aw9d9j29j2d9aj9d!23", 10000000, nil, false))

	route.Kid("/auth", middleware.MemberAuth(), controller2.Router())

	route.Kid("/", middleware.MemberAuth(), controller.Router())

	route.Kid("/mybabies", middleware.MemberAuth(), controller3.Router())
	route.Kid("/babyGrowning", middleware.MemberAuth(), daily.Router())

	return route
}
