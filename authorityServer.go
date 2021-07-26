package authorityServer // import "github.com/hypwxm/authorityServer"

import (
	"authorityServer/config"
	"authorityServer/router"

	"github.com/hypwxm/rider"
)

func Init(app rider.Rider, cfg config.ConfigDefine) *rider.Router {
	config.InitConfig(cfg)
	// app.Kid("/", router.Router())
	return router.Router()
}
