package config

import "net"

type ListenerType int32

const (
	LISTENER_STATIC ListenerType = 0
	LISTENER_EDS    ListenerType = 3
)

type FilterChainMatch_ConnectionSourceType int32

const (
	// Any connection source matches.
	FilterChainMatch_ANY FilterChainMatch_ConnectionSourceType = 0
	// Match a connection originating from the same host.
	FilterChainMatch_SAME_IP_OR_LOOPBACK FilterChainMatch_ConnectionSourceType = 1
	// Match a connection originating from a different host.
	FilterChainMatch_EXTERNAL FilterChainMatch_ConnectionSourceType = 2
)

// Listener
type Listener interface {
	GetName() string
	GetAddress() net.Addr
	GetListenerFilters() []Filter
	GetFilterChains() []FilterChain
	GetDefaultFilterChain() FilterChain
	GetUseOriginalDst() bool
	GetBindToPort() bool
}

type FilterChain interface {
	GetName() string
	GetFilterChainMatch() FilterChainMatch
	GetFilters() []Filter
	SetName(string)
}

type FilterChainMatch interface {
	GetPrefixRanges() []CidrRange
	GetServerNames() []string
	GetApplicationProtocols() []string
	GetDirectSourcePrefixRanges() []CidrRange
	GetSourceType() FilterChainMatch_ConnectionSourceType
	GetSourcePrefixRanges() []CidrRange
	GetSourcePorts() []uint32
	GetDestinationPort() uint32
	GetTransportProtocol() string
}
