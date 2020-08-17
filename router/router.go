package router

import (
	"github.com/hypwxm/rider"
	"babygrowing/middleware"
	controller8 "babygrowing/service/admin/role/controller"
	controller13 "babygrowing/service/admin/rolePermission/menu/controller"
	controller15 "babygrowing/service/admin/rolePermission/source/controller"
	controller7 "babygrowing/service/admin/user/controller"
	controller2 "babygrowing/service/auth/controller"
	controller4 "babygrowing/service/house/controller"
	controller6 "babygrowing/service/matter/matter/controller"
	controller9 "babygrowing/service/matter/matterElement/controller"
	controller10 "babygrowing/service/matter/matterElementOption/controller"
	controller11 "babygrowing/service/matter/matterVisible/controller"
	controller12 "babygrowing/service/menu/controller"
	controller5 "babygrowing/service/newsDynamics/controller"
	controller16 "babygrowing/service/societies/base/controller"
	"babygrowing/service/user/controller"
	controller3 "babygrowing/service/vote/controller"
	controller14 "babygrowing/service/webSource/controller"
)

func Router() *rider.Router {
	route := rider.NewRouter()
	route.USE(rider.RiderJwt("Authorization", "ni2QWN29DJQJDNI923N=-230S-23!23", 10000000, nil, false))
	route.Kid("/user", middleware.Auth(), controller.Router())
	route.Kid("/auth", middleware.Auth(), controller2.Router())
	route.Kid("/vote", middleware.Auth(), controller3.Router())
	route.Kid("/house", middleware.Auth(), controller4.Router())
	route.Kid("/newsDynamics", middleware.Auth(), controller5.Router())
	route.Kid("/matter", middleware.Auth(), controller6.Router())
	route.Kid("/matterElement", middleware.Auth(), controller9.Router())
	route.Kid("/matterElementOption", middleware.Auth(), controller10.Router())
	route.Kid("/matterVisible", middleware.Auth(), controller11.Router())
	route.Kid("/settings/menu", middleware.Auth(), controller12.Router())
	route.Kid("/settings/webSource", middleware.Auth(), controller14.Router())
	route.Kid("/societies", middleware.Auth(), controller16.Router())

	route.Kid("/adminuser", middleware.Auth(), controller7.Router())
	route.Kid("/adminrole", middleware.Auth(), controller8.Router())
	route.Kid("/adminrole/menuPermission", middleware.Auth(), controller13.Router())
	route.Kid("/adminrole/sourcePermission", middleware.Auth(), controller15.Router())

	return route
}
