package config

import "net"

type LocalityLbEndpointsImpl struct {
	LbEndpoints []*LbEndpointImpl `json:"lb_endpoints"`
}

func (e *LocalityLbEndpointsImpl) GetLbEndpoint() []LbEndpoint {
	var edps []LbEndpoint
	for _, edp := range e.LbEndpoints {
		edps = append(edps, edp)
	}
	return edps
}

type LbEndpointImpl struct {
	Endpoint            *EndpointImpl `json:"endpoint"`
	LoadBalancingWeight int           `json:"load_balancing_weight"`
}

func (e *LbEndpointImpl) GetEndpoint() Endpoint {
	return e.Endpoint
}

func (e *LbEndpointImpl) GetLoadBalancingWeight() int {
	return e.LoadBalancingWeight
}

func (e *LbEndpointImpl) SetLoadBalancingWeight(w int) {
	e.LoadBalancingWeight = w
}

type EndpointImpl struct {
	Address *Address `json:"address"`
}

func (e *EndpointImpl) GetAddress() net.Addr {
	return e.Address
}
