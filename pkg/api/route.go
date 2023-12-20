package api

import (
	"github.com/wereliang/sota-mesh/pkg/config"
)

// RouteConfigMatcher
type RouteConfigMatcher interface {
	Match(RequestHeader) config.RouteEntry
}

// RouteConfigManager
type RouteConfigManager interface {
	GetRouteConfig(string) config.RouteConfiguration
	AddOrUpdateRouteConfig(config.RouteConfigType, config.RouteConfiguration) error
	DeleteRouteConfig(string) error
}
