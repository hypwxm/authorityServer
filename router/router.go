package router

import (
	"github.com/hypwxm/rider"
)

func Router() *rider.Router {
	route := rider.NewRouter()
	route.Kid("/", AdminRouter())
	return route
}
