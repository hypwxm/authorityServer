package baseController

import "github.com/hypwxm/rider"

func Router() *rider.Router {
	router := rider.NewRouter()
	router.POST("/upload", upload)
	return router
}
