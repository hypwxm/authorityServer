package router

import (
	"babygrowing/middleware"
	controller2 "babygrowing/service/auth/controller"
	daily "babygrowing/service/member/daily/controller"
	controller3 "babygrowing/service/member/mybabies/controller"
	"babygrowing/service/member/user/controller"

	"github.com/hypwxm/rider"
)

func AppRouter() *rider.Router {
	route := rider.NewRouter()
	route.USE(rider.RiderJwt("AppAuth", "ni2QWN2asd23aw9d9j29j2d9aj9d!23", 10000000, nil, false))

	route.Kid("/auth", middleware.AppAuth(), controller2.Router())

	route.Kid("/user", middleware.Auth(), controller.Router())

	route.Kid("/mybabies", middleware.Auth(), controller3.Router())
	route.Kid("/babyGrowning", middleware.Auth(), daily.Router())

	return route
}
