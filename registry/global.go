package registry

//var eventRegistry = &EventRegistry{}

func Instance() *EventRegistry {
	return eventRegistry
}

func init() {
	eventRegistry.RegisterPresetMainnet()
}
