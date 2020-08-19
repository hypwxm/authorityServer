package router

import (
	"babygrowing/middleware"
	controller2 "babygrowing/service/auth/controller"
	appController2 "babygrowing/service/user/appController"

	"github.com/hypwxm/rider"
)

func AppRouter() *rider.Router {
	route := rider.NewRouter()
	route.USE(rider.RiderJwt("AppAuth", "ni2QWN2asd23aw9d9j29j2d9aj9d!23", 10000000, nil, false))

	route.Kid("/auth", middleware.AppAuth(), controller2.Router())

	route.Kid("/user", middleware.AppAuth(), appController2.Router())

	return route
}
