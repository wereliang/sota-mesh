package config

// Bootstrap
type Bootstrap interface {
	GetStaticResources() StaticResources
}

// StaticResources Bootstrap's static resources
type StaticResources interface {
	GetListeners() []Listener
	GetClusters() []Cluster
}
