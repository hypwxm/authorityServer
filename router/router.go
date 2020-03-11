package router

import (
	"github.com/hypwxm/rider"
	"worldbar/middleware"
	controller2 "worldbar/service/auth/controller"
	controller4 "worldbar/service/house/controller"
	"worldbar/service/user/controller"
	controller3 "worldbar/service/vote/controller"
)

func Router() *rider.Router {
	route := rider.NewRouter()
	route.USE(rider.RiderJwt("Authorization", "ni2QWN29DJQJDNI923N=-230S-23!23", 10000000, nil))
	route.Kid("/user", middleware.Auth(), controller.Router())
	route.Kid("/auth", middleware.Auth(), controller2.Router())
	route.Kid("/vote", middleware.Auth(), controller3.Router())
	route.Kid("/house", middleware.Auth(), controller4.Router())
	return route
}
