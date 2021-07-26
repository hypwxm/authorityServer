package controller

import (
	"github.com/hypwxm/rider"
)

func Router() *rider.Router {
	route := rider.NewRouter()

	route.POST("/open/login", login)
	route.POST("/loginUser", loginAdmin)

	return route
}
