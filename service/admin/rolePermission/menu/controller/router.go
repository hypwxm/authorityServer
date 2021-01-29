package controller

import "github.com/hypwxm/rider"

func Router() *rider.Router {
	route := rider.NewRouter()

	route.POST("/create", create)
	route.POST("/list", list)
	route.POST("/del", del)

	return route
}
