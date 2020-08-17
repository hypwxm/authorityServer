package baseControl

import (
	"babygrowing/base_control/controller"

	"github.com/hypwxm/rider"
)

func Init(app *rider.Rider) {
	app.Kid("/server/common", controller.Router())
}
