package api

import (
	"net"

	"github.com/wereliang/sota-mesh/pkg/config"
)

// ActiveConnection active connection wrapper
type ActiveConnection interface {
	FilterManager
	Connection
	OnLoop()
}

// ActiveListener is listener handler
type ActiveListener interface {
	SotaObject
	ListenerFilterManager

	Start() error
	Type() config.ListenerType
	Listener() Listener
	GetUseOriginalDst() bool
	GetBindToPort() bool
}

// ListenerManager
type ListenerManager interface {
	AddOrUpdateListener(config.ListenerType, config.Listener) error
	FindListenerByAddress(net.Addr) ActiveListener
	FindListenerByName(string) ActiveListener
	DeleteListener(string) error
}
