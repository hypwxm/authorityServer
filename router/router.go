package router

import (
	"github.com/hypwxm/rider"
	"worldbar/middleware"
	controller8 "worldbar/service/admin/role/controller"
	controller7 "worldbar/service/admin/user/controller"
	controller2 "worldbar/service/auth/controller"
	controller4 "worldbar/service/house/controller"
	controller6 "worldbar/service/matter/matter/controller"
	controller5 "worldbar/service/newsDynamics/controller"
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
	route.Kid("/newsDynamics", middleware.Auth(), controller5.Router())
	route.Kid("/matter", middleware.Auth(), controller6.Router())
	route.Kid("/adminuser", middleware.Auth(), controller7.Router())
	route.Kid("/adminrole", middleware.Auth(), controller8.Router())

	return route
}
