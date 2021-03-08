package event

import (
	evnbus "github.com/asaskevich/EventBus"
)

var Ebus evnbus.Bus

func init() {
	Ebus = evnbus.New()
}
