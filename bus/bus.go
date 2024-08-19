package bus

import (
	evbus "github.com/asaskevich/EventBus"
)

var Bus evbus.Bus

func init() {
	Bus = evbus.New()
}
