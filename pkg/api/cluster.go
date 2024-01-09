package api

import "github.com/wereliang/sota-mesh/pkg/config"

// Cluster handler upstream endpoints
type Cluster interface {
	SotaObject
	Snapshot() ClusterSnapshot
	UpdateEndpoints(config.LbEndpointSet)
	GetQosRouter() QosRouter
	Close()
}

// ClusterSnapshot is a thread-safe cluster snapshot
type ClusterSnapshot interface {
	EndpointSet() config.LbEndpointSet
	ClusterInfo() config.Cluster
	// LoadBalancer() LoadBalancer
}

// ClusterManager is a manager for cluster
type ClusterManager interface {
	AddOrUpdateCluster(config.Cluster) error
	DeleteCluster(string) error
	GetCluster(string) Cluster
	UpdateClusterEndpoints(string, config.LbEndpointSet) error
}
