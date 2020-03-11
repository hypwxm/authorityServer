package controller

import (
	"worldbar/base_control/controller/baseController"

	"github.com/hypwxm/rider"
)

func Router() *rider.Router {
	router := rider.NewRouter()
	router.Kid("/base", baseController.Router())
	return router
}
