package appController

import "github.com/hypwxm/rider"

func Router() *rider.Router {
	route := rider.NewRouter()

	//
	route.POST("/list", list)
	route.POST("/create", create)
	route.POST("/del", del)

	return route
}
