package api

// ListenerFilter
type ListenerFilter interface {
	OnAccept(ListenerFilterCallbacks) FilterStatus
}

// ListenerFilterCallbacks
type ListenerFilterCallbacks interface {
	Connection() Connection
}

// ListenerFilterManager manager listener filter
type ListenerFilterManager interface {
	AddAcceptFilter(ListenerFilter)
}

type ListenerFilterCreator func(ListenerFilterManager)

// ListernerFactory for listener factory
type ListernerFactory interface {
	Factory
	CreateFilterFactory(interface{}, FactoryContext) ListenerFilterCreator
}
