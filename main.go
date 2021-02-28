package main

import (
	"babygrow/config"
	"babygrow/router"
	"log"
	"os"
	"path/filepath"

	"github.com/hypwxm/rider"
	"github.com/hypwxm/rider/modules"
)

func main() {

	app := rider.New()
	app.Logger(8)

	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)

	wd, _ := os.Getwd()
	app.SetStatic(filepath.Join(wd, "assets"), "/assets")

	app.Kid("/", router.Router())
	// app.Kid("/app", router.AppRouter())

	// modules.DefaultSecureConfig.XFrameOptions = "ALLOW-FROM http://localhost:9527"
	app.USE(modules.SecureHeader())
	app.Graceful(config.Config.ServerPort)

}
