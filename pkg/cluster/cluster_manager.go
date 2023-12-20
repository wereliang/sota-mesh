package cluster

import (
	"fmt"
	"sync"

	"github.com/wereliang/sota-mesh/pkg/api"
	"github.com/wereliang/sota-mesh/pkg/config"
	"github.com/wereliang/sota-mesh/pkg/log"
)

func NewClusterManager(clusters []config.Cluster) (api.ClusterManager, error) {
	cm := &clusterManager{}
	for _, cluster := range clusters {
		err := cm.AddOrUpdateCluster(cluster)
		if err != nil {
			log.Error("add cluster error:%s", err)
		}
	}
	return cm, nil
}

type clusterManager struct {
	clusterMap sync.Map
}

func (cm *clusterManager) AddOrUpdateCluster(c config.Cluster) error {

	if cluster := cm.GetCluster(c.GetName()); cluster != nil {
		log.Debug("cluster %s close", c.GetName())
		cluster.Close()
	}

	newCluster, err := NewCluster(c)
	if err != nil {
		return err
	}
	cm.clusterMap.Store(c.GetName(), newCluster)
	return nil
}

func (cm *clusterManager) GetCluster(name string) api.Cluster {
	if c, ok := cm.clusterMap.Load(name); ok {
		return c.(api.Cluster)
	}
	return nil
}

func (cm *clusterManager) DeleteCluster(name string) error {
	cm.clusterMap.Delete(name)
	return nil
}

func (cm *clusterManager) UpdateClusterEndpoints(clusterName string, edps config.LbEndpointSet) error {
	cluster := cm.GetCluster(clusterName)
	if cluster == nil {
		return fmt.Errorf("update cluster hosts fail. not found cluster:%s", clusterName)
	}
	cluster.UpdateEndpoints(edps)
	return nil
}
