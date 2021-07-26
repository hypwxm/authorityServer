package baseController

import (
	"github.com/hypwxm/authorityServer/middleware"

	"github.com/hypwxm/rider"
)

func Router() *rider.Router {
	router := rider.NewRouter()
	router.POST("/upload", middleware.Auth(), uploadToAliyunOss)

	router.POST("/customUpload", middleware.Auth(), uploadToAliyunOss)
	return router
}
