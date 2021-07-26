package registry

import "github.com/umbracle/go-web3"

var eventRegistry = NewEventRegistry()

func Instance() *EventRegistry {
	return eventRegistry
}

func init() {
	eventRegistry.RegisterPresetMainnet()
	web3.RegisterParser(eventRegistry)
}
