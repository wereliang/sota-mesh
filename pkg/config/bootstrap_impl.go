package config

import (
	"github.com/mcuadros/go-defaults"
)

type BootstrapImpl struct {
	StaticResources *StaticResourcesImpl `json:"static_resources"`
}

func (b *BootstrapImpl) GetStaticResources() StaticResources {
	return b.StaticResources
}

type StaticResourcesImpl struct {
	Listeners []*ListenerImpl `json:"listeners"`
	Clusters  []*ClusterImpl  `json:"clusters"`
}

func (s *StaticResourcesImpl) GetListeners() []Listener {
	var listeners []Listener
	for _, l := range s.Listeners {
		// set default value, by "default" tag
		defaults.SetDefaults(l)
		listeners = append(listeners, l)
	}
	return listeners
}

func (s *StaticResourcesImpl) GetClusters() []Cluster {
	var clusters []Cluster
	for _, c := range s.Clusters {
		clusters = append(clusters, c)
	}
	return clusters
}
