package authorityServer // import https://gitee.com/hypwxm_admin/authority-server

import (
	"authorityServer/config"
	"authorityServer/router"

	_ "authorityServer/event"
	_ "authorityServer/service/message/model"

	"github.com/hypwxm/rider"
)

func Init(app rider.Rider, cfg config.ConfigDefine) {
	config.InitConfig(cfg)
	app.Kid("/", router.Router())
}
