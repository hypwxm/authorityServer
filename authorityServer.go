package authorityServer // import "github.com/hypwxm/authorityServer"

import (
	"github.com/hypwxm/authorityServer/config"
	"github.com/hypwxm/authorityServer/router"

	"github.com/hypwxm/rider"
)

func Init(app *rider.Rider, cfg config.ConfigDefine) *rider.Router {
	config.InitConfig(cfg)
	// app.Kid("/", router.Router())
	return router.Router()
}
