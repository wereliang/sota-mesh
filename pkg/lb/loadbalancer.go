package lb

import (
	"github.com/wereliang/sota-mesh/pkg/api"
	"github.com/wereliang/sota-mesh/pkg/config"
)

type LBCreator func(config.LbEndpointSet) api.LoadBalancer

var (
	lbFactory = make(map[config.LoadBalancerType]LBCreator)
)

func registLoadBalancer(t config.LoadBalancerType, creator LBCreator) {
	lbFactory[t] = creator
}

func NewLoadBalancer(t config.LoadBalancerType, edps config.LbEndpointSet) api.LoadBalancer {
	if creator, ok := lbFactory[t]; ok {
		return creator(edps)
	}
	return nil
}
