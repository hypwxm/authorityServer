package main

import (
	"babygrow/config"
	"babygrow/router"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "babygrow/event"
	_ "babygrow/service/message/model"

	"github.com/hypwxm/rider"
	"github.com/hypwxm/rider/modules"
)

func main() {

	app := rider.New()
	app.SetHttpReadTimeout(120 * time.Second)
	app.Logger(8)

	if config.Env == config.ENV_DEV {
		// 开发环境下打开跨域，允许swagger等接口在线请求
		app.SetAccessCtl(func(c rider.Context) *rider.AccessControl {
			return &rider.AccessControl{
				AccessControlAllowCredentials: "true",
				AccessControlAllowOrigin:      "*",
				AccessControlAllowHeaders:     "*",
			}
		})
	}

	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)

	wd, _ := os.Getwd()
	app.SetStatic(filepath.Join(wd, "assets"), "/assets")

	app.Kid("/", router.Router())
	// app.Kid("/app", router.AppRouter())

	// modules.DefaultSecureConfig.XFrameOptions = "ALLOW-FROM http://localhost:9527"
	app.USE(modules.SecureHeader())
	app.Listen(config.Config.ServerPort)

}
