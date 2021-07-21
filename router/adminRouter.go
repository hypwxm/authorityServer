package router

import (
	"authorityServer/middleware"
	menuController "authorityServer/service/admin/menu/controller"
	orgController "authorityServer/service/admin/org/controller"
	roleMenuController "authorityServer/service/admin/rolePermission/menu/controller"

	controller8 "authorityServer/service/admin/role/controller"
	controller7 "authorityServer/service/admin/user/controller"

	controller2 "authorityServer/service/auth/controller"

	baseControl "authorityServer/base_control"
	controller12 "authorityServer/service/menu/controller"
	controller14 "authorityServer/service/webSource/controller"

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
