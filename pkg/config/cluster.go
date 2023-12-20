package config

type ClusterType int32

const (
	Cluster_Static       ClusterType = 0
	Cluster_Strict_DNS   ClusterType = 1
	Cluster_Logical_DNS  ClusterType = 2
	Cluster_EDS          ClusterType = 3
	Cluster_ORIGINAL_DST ClusterType = 4
)

// ClusterInfo defines a cluster's information
type Cluster interface {
	GetName() string
	GetType() ClusterType
	GetLbPolicy() LoadBalancerType
	SetLbPolicy(LoadBalancerType)
	GetLoadAssignment() ClusterLoadAssignment
	GetUpstreamBindConfig() BindConfig
}

// ClusterLoadAssignment
type ClusterLoadAssignment interface {
	GetClusterName() string
	GetEndpoints() []LocalityLbEndpoints
}
