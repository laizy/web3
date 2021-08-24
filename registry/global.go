package registry

import "github.com/laizy/web3"

var eventRegistry = NewEventRegistry()

func Instance() *EventRegistry {
	return eventRegistry
}

func init() {
	eventRegistry.RegisterPresetMainnet()
	web3.RegisterParser(eventRegistry)
}
