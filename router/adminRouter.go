package router

import (
	"github.com/hypwxm/authorityServer/middleware"
	menuController "github.com/hypwxm/authorityServer/service/admin/menu/controller"
	orgController "github.com/hypwxm/authorityServer/service/admin/org/controller"
	roleMenuController "github.com/hypwxm/authorityServer/service/admin/rolePermission/menu/controller"

	controller8 "github.com/hypwxm/authorityServer/service/admin/role/controller"
	controller7 "github.com/hypwxm/authorityServer/service/admin/user/controller"

	controller2 "github.com/hypwxm/authorityServer/service/auth/controller"

	baseControl "github.com/hypwxm/authorityServer/base_control"

	"github.com/hypwxm/rider"
)

func AdminRouter() *rider.Router {
	route := rider.NewRouter()
	route.USE(rider.RiderJwt("Authorization", "ni2QWN29DJQJDNI923N=-230S-23!23", 10000000, nil, false))

	// 基础模块，图片上传等
	baseControl.Init(route)

	route.Kid("/auth", middleware.Auth(), controller2.Router())

	route.Kid("/adminuser", middleware.Auth(), controller7.Router())
	route.Kid("/adminrole", middleware.Auth(), controller8.Router())
	route.Kid("/adminorg", middleware.Auth(), orgController.Router())
	route.Kid("/adminmenu", middleware.Auth(), menuController.Router())
	route.Kid("/adminrolemenu", middleware.Auth(), roleMenuController.Router())

	return route
}
