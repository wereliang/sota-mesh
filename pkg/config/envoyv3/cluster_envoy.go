package envoyv3

import (
	"net"

	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_config_endpoint_v3 "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	"github.com/wereliang/sota-mesh/pkg/config"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

type ClusterEnvoy struct {
	*envoy_config_cluster_v3.Cluster
}

func (c *ClusterEnvoy) GetType() config.ClusterType {
	return config.ClusterType(c.Cluster.GetType())
}

func (c *ClusterEnvoy) GetLbPolicy() config.LoadBalancerType {
	return config.LoadBalancerType(c.Cluster.GetLbPolicy())
}

func (c *ClusterEnvoy) SetLbPolicy(lbt config.LoadBalancerType) {
	c.Cluster.LbPolicy = envoy_config_cluster_v3.Cluster_LbPolicy(lbt)
}

func (c *ClusterEnvoy) GetLoadAssignment() config.ClusterLoadAssignment {
	if x := c.Cluster.GetLoadAssignment(); x != nil {
		return &ClusterLoadAssignmentEnvoy{x}
	}
	return nil
}

func (c *ClusterEnvoy) GetUpstreamBindConfig() config.BindConfig {
	if x := c.Cluster.GetUpstreamBindConfig(); x != nil {
		return &BindConfigEnvoy{x}
	}
	return nil
}

type BindConfigEnvoy struct {
	*envoy_config_core_v3.BindConfig
}

func (b *BindConfigEnvoy) GetSourceAddress() *config.SocketAddress {
	addr := b.BindConfig.GetSourceAddress()
	return &config.SocketAddress{
		Protocol:  addr.GetProtocol().String(),
		Address:   addr.GetAddress(),
		PortValue: addr.GetPortValue(),
	}
}

type ClusterLoadAssignmentEnvoy struct {
	*envoy_config_endpoint_v3.ClusterLoadAssignment
}

func (c *ClusterLoadAssignmentEnvoy) GetEndpoints() []config.LocalityLbEndpoints {
	var edps []config.LocalityLbEndpoints
	for _, edp := range c.ClusterLoadAssignment.GetEndpoints() {
		edps = append(edps, &LocalityLbEndpointsEnvoy{edp})
	}
	return edps
}

type LocalityLbEndpointsEnvoy struct {
	*envoy_config_endpoint_v3.LocalityLbEndpoints
}

func (e *LocalityLbEndpointsEnvoy) GetLbEndpoint() []config.LbEndpoint {
	var edps []config.LbEndpoint
	for _, edp := range e.LocalityLbEndpoints.GetLbEndpoints() {
		edps = append(edps, &LbEndpointEnvoy{edp})
	}
	return edps
}

type LbEndpointEnvoy struct {
	*envoy_config_endpoint_v3.LbEndpoint
}

func (e *LbEndpointEnvoy) GetEndpoint() config.Endpoint {
	return &EndpointEnvoy{e.LbEndpoint.GetEndpoint()}
}

func (e *LbEndpointEnvoy) GetLoadBalancingWeight() int {
	if e.LbEndpoint.GetLoadBalancingWeight() != nil {
		return int(e.LbEndpoint.GetLoadBalancingWeight().GetValue())
	}
	return 0
}

func (e *LbEndpointEnvoy) SetLoadBalancingWeight(w int) {
	e.LbEndpoint.LoadBalancingWeight = wrapperspb.UInt32(uint32(w))
}

type EndpointEnvoy struct {
	*envoy_config_endpoint_v3.Endpoint
}

func (e *EndpointEnvoy) GetAddress() net.Addr {
	addr, err := ToNetAddr(e.Endpoint.GetAddress())
	if err != nil {
		panic(err)
	}
	return addr
}
