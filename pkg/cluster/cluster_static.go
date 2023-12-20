package cluster

import (
	"github.com/wereliang/sota-mesh/pkg/api"
	"github.com/wereliang/sota-mesh/pkg/config"
)

type staticCluster struct {
	*SimpleCluster
}

func newStaticCluster(c config.Cluster) (api.Cluster, error) {
	simple := newSimpleCluster(c)
	edps, _ := simple.getConfigEndpoints()
	simple.UpdateEndpoints(edps)
	static := &staticCluster{simple}
	return static, nil
}

func init() {
	registCluster(config.Cluster_Static,
		func(c config.Cluster) (api.Cluster, error) {
			return newStaticCluster(c)
		})
}
