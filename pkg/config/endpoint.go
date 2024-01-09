package config

import (
	"net"
	"time"
)

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
	GetExpired() time.Time // 过载过期时间
	SetExpired(t time.Time)
}

type Endpoint interface {
	GetAddress() net.Addr
}
