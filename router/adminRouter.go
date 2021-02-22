package router

import (
	"babygrowing/middleware"
	menuController "babygrowing/service/admin/menu/controller"
	orgController "babygrowing/service/admin/org/controller"
	roleMenuController "babygrowing/service/admin/rolePermission/menu/controller"

	controller8 "babygrowing/service/admin/role/controller"
	controller7 "babygrowing/service/admin/user/controller"

	controller2 "babygrowing/service/auth/controller"

	baseControl "babygrowing/base_control"
	controller12 "babygrowing/service/menu/controller"
	controller14 "babygrowing/service/webSource/controller"

	"github.com/hypwxm/rider"
)

func AdminRouter() *rider.Router {
	route := rider.NewRouter()
	route.USE(rider.RiderJwt("Authorization", "ni2QWN29DJQJDNI923N=-230S-23!23", 10000000, nil, false))

	// 基础模块，图片上传等
	baseControl.Init(route)

	route.Kid("/auth", middleware.Auth(), controller2.Router())

	route.Kid("/settings/menu", middleware.Auth(), controller12.Router())
	route.Kid("/settings/webSource", middleware.Auth(), controller14.Router())

	route.Kid("/adminuser", middleware.Auth(), controller7.Router())
	route.Kid("/adminrole", middleware.Auth(), controller8.Router())
	route.Kid("/adminorg", middleware.Auth(), orgController.Router())
	route.Kid("/adminmenu", middleware.Auth(), menuController.Router())
	route.Kid("/adminrolemenu", middleware.Auth(), roleMenuController.Router())

	return route
}