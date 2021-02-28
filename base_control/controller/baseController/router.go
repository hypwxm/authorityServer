package baseController

import (
	"babygrow/middleware"

	"github.com/hypwxm/rider"
)

func Router() *rider.Router {
	router := rider.NewRouter()
	router.POST("/upload", middleware.Auth(), uploadToAliyunOss)

	router.POST("/customUpload", middleware.MemberAuth(), uploadToAliyunOss)
	return router
}
