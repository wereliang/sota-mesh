package api

import (
	"github.com/wereliang/sota-mesh/pkg/config"
)

type FilterStatus int

const (
	Stop     FilterStatus = 0
	Continue FilterStatus = 1
)

// Factory is basic factory
type Factory interface {
	Name() string
	CreateEmptyConfigProto() interface{}
}

// FactoryContext some context for filter
type FactoryContext interface {
	ListenerManager() ListenerManager
	ClusterManager() ClusterManager
	RouteConfigManager() RouteConfigManager
}

// FilterChainManager filter chain manager
type FilterChainManager interface {
	AddFilterChains([]config.FilterChain, config.FilterChain)
	FindFilterChains(ConnectionContext) config.FilterChain
}
