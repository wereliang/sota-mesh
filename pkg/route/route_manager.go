package route

import (
	"sync"

	"github.com/wereliang/sota-mesh/pkg/config"
)

type RouteManagerImpl struct {
	routeConfigMap sync.Map
}

func (rm *RouteManagerImpl) GetRouteConfig(name string) config.RouteConfiguration {
	if r, ok := rm.routeConfigMap.Load(name); ok {
		return r.(config.RouteConfiguration)
	}
	return nil
}

func (rm *RouteManagerImpl) AddOrUpdateRouteConfig(typ config.RouteConfigType, rc config.RouteConfiguration) error {
	rm.routeConfigMap.Store(rc.GetName(), rc)
	return nil
}

func (rm *RouteManagerImpl) DeleteRouteConfig(name string) error {
	rm.routeConfigMap.Delete(name)
	return nil
}
