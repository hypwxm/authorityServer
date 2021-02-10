package controller

import (
	"github.com/hypwxm/rider"
)

func Router() *rider.Router {
	route := rider.NewRouter()

	route.POST("/open/login", memberLogin)
	route.POST("/loginUser", loginAppUser)

	return route
}
