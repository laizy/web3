package registry

import (
	"github.com/laizy/web3"
)

var eventRegistry = NewEventRegistry()
var errorRegister = NewErrorRegistry()

func Instance() *EventRegistry {
	return eventRegistry
}

func ErrInstance() *ErrorRegistry {
	return errorRegister
}

func init() {
	eventRegistry.RegisterPresetMainnet()
	web3.RegisterParser(eventRegistry)
}
