package baseControl

import (
	"babygrow/base_control/controller"

	"github.com/hypwxm/rider"
)

func Init(app *rider.Router) {
	app.Kid("/common", controller.Router())
}
