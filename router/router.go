package router

import (
	"babygrowing/middleware"
	controller8 "babygrowing/service/admin/role/controller"
	controller13 "babygrowing/service/admin/rolePermission/menu/controller"
	controller15 "babygrowing/service/admin/rolePermission/source/controller"
	controller7 "babygrowing/service/admin/user/controller"
	controller2 "babygrowing/service/auth/controller"
	controller12 "babygrowing/service/menu/controller"
	controller3 "babygrowing/service/mybabies/controller"
	"babygrowing/service/user/controller"
	controller14 "babygrowing/service/webSource/controller"
	"github.com/hypwxm/rider"
)

func Router() *rider.Router {
	route := rider.NewRouter()
	route.USE(rider.RiderJwt("Authorization", "ni2QWN29DJQJDNI923N=-230S-23!23", 10000000, nil, false))
	route.Kid("/user", middleware.Auth(), controller.Router())
	route.Kid("/auth", middleware.Auth(), controller2.Router())

	route.Kid("/settings/menu", middleware.Auth(), controller12.Router())
	route.Kid("/settings/webSource", middleware.Auth(), controller14.Router())
	route.Kid("/mybabies", middleware.Auth(), controller3.Router())

	route.Kid("/adminuser", middleware.Auth(), controller7.Router())
	route.Kid("/adminrole", middleware.Auth(), controller8.Router())
	route.Kid("/adminrole/menuPermission", middleware.Auth(), controller13.Router())
	route.Kid("/adminrole/sourcePermission", middleware.Auth(), controller15.Router())

	return route
}
