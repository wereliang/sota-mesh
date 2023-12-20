package config

type RouteConfigType int32

const (
	ROUTE_CONFIG_STATIC RouteConfigType = 0
	ROUTE_CONFIG_EDS    RouteConfigType = 3
)

type RouteConfiguration interface {
	GetName() string
	GetVirtualHosts() []VirtualHost
}

type RouteEntry interface {
	GetClusterName() string
}

type VirtualHost interface {
	GetName() string
	GetDomains() []string
	GetRoutes() []Route
}

type Route interface {
	GetMatch() RouteMatch
	GetRoute() RouteAction
}

type RouteMatch interface {
	GetPath() string
	GetPrefix() string
}

type RouteAction interface {
	GetCluster() string
}
