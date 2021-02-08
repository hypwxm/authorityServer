package controller

import (
	"github.com/hypwxm/rider"
)

func Router() *rider.Router {
	route := rider.NewRouter()

	route.POST("/open/admin/login", login)
	route.POST("/admin/loginUser", loginAdmin)

	route.POST("/open/member/login", memberLogin)
	route.POST("/app/loginUser", loginAppUser)

	return route
}
