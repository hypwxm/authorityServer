package router

import (
	"worldbar/middleware"
	controller2 "worldbar/service/auth/controller"
	"worldbar/service/house/appController"
	appController2 "worldbar/service/user/appController"

	"github.com/hypwxm/rider"
)

func AppRouter() *rider.Router {
	route := rider.NewRouter()
	route.USE(rider.RiderJwt("AppAuth", "ni2QWN2asd23aw9d9j29j2d9aj9d!23", 10000000, nil))

	route.Kid("/auth", middleware.AppAuth(), controller2.Router())

	route.Kid("/user", middleware.AppAuth(), appController2.Router())
	route.Kid("/house", middleware.AppAuth(), appController.Router())


	return route
}
