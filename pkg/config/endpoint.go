package config

import "net"

const (
	Default_LoadBalance_Weight = 100
)

type LocalityLbEndpoints interface {
	GetLbEndpoint() []LbEndpoint
}

type LbEndpointSet []LbEndpoint

type LbEndpoint interface {
	GetEndpoint() Endpoint
	GetLoadBalancingWeight() int
	SetLoadBalancingWeight(int)
}

type Endpoint interface {
	GetAddress() net.Addr
}
