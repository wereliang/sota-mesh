package envoyv3

import (
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/wereliang/sota-mesh/pkg/config"
)

type RouteConfigurationEnvoy struct {
	*envoy_config_route_v3.RouteConfiguration
}

func (rc *RouteConfigurationEnvoy) GetVirtualHosts() []config.VirtualHost {
	var vhs []config.VirtualHost
	for _, vh := range rc.RouteConfiguration.GetVirtualHosts() {
		vhs = append(vhs, &VirtualHostEnvoy{vh})
	}
	return vhs
}

type VirtualHostEnvoy struct {
	*envoy_config_route_v3.VirtualHost
}

func (vh *VirtualHostEnvoy) GetRoutes() []config.Route {
	var routes []config.Route
	for _, r := range vh.VirtualHost.GetRoutes() {
		routes = append(routes, &RouteEnvoy{r})
	}
	return routes
}

type RouteEnvoy struct {
	*envoy_config_route_v3.Route
}

func (r *RouteEnvoy) GetMatch() config.RouteMatch {
	if x := r.Route.GetMatch(); x != nil {
		return &RouteMatchEnvoy{x}
	}
	return nil
}

func (r *RouteEnvoy) GetRoute() config.RouteAction {
	return r.Route.GetRoute()
}

type RouteMatchEnvoy struct {
	*envoy_config_route_v3.RouteMatch
}
