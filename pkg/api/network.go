package api

import (
	"net"
)

// Listener network listener wrapper
type Listener interface {
	Addr() net.Addr
	Listen() error
	SetCallback(ListenerCallback)
}

// ListenerCallback
type ListenerCallback interface {
	OnAccept(conn Connection)
}

// Connection network connection wrapper
type Connection interface {
	net.Conn
	Raw() net.Conn
	Context() ConnectionContext
	Peek(n int) ([]byte, error)
}

// ConnectionCallbacks
type ConnectionCallbacks interface {
	Connection
}

// ConnectionContext
type ConnectionContext interface {
	GetDestinationPort() uint32
	GetDestinationIP() net.IP
	GetServerName() string
	GetTransportProtocol() string
	GetApplicationProtocol() string
	GetDirectSourceIP() net.IP
	GetSourceType() int32
	GetSourceIP() net.IP
	GetSourcePort() uint32
	LocalAddressRestored() bool
	SetOriginalDestination(net.IP, uint32)
	SetApplicationProtocol(string)
}
