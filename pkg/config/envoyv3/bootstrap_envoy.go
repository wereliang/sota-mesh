package envoyv3

import (
	envoy_config_bootstrap_v3 "github.com/envoyproxy/go-control-plane/envoy/config/bootstrap/v3"
	"github.com/wereliang/sota-mesh/pkg/config"
)

type BootstrapEnvoy struct {
	*envoy_config_bootstrap_v3.Bootstrap
}

func (b *BootstrapEnvoy) GetStaticResources() config.StaticResources {
	if x := b.Bootstrap.GetStaticResources(); x != nil {
		return &StaticResourcesEnvoy{x}
	}
	return nil
}

type StaticResourcesEnvoy struct {
	staticResources *envoy_config_bootstrap_v3.Bootstrap_StaticResources
}

func (s *StaticResourcesEnvoy) GetListeners() []config.Listener {
	var listeners []config.Listener
	for _, l := range s.staticResources.GetListeners() {
		listeners = append(listeners, &ListenerEnvoy{l})
	}
	return listeners
}

func (s *StaticResourcesEnvoy) GetClusters() []config.Cluster {
	var clusters []config.Cluster
	for _, c := range s.staticResources.GetClusters() {
		clusters = append(clusters, &ClusterEnvoy{c})
	}
	return clusters
}
