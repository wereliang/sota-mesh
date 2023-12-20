package api

import "github.com/wereliang/sota-mesh/pkg/config"

// LoadBalancerContext use for lb context
type LoadBalancerContext interface {
	Connection() Connection
}

// LoadBalancer
type LoadBalancer interface {
	// Select return host by special algorithm
	Select(LoadBalancerContext) config.LbEndpoint
}
