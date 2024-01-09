package cluster

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/wereliang/sota-mesh/pkg/api"
	"github.com/wereliang/sota-mesh/pkg/config"
	"github.com/wereliang/sota-mesh/pkg/log"
	"github.com/wereliang/sota-mesh/pkg/qos"
)

type ClusterCreator func(config.Cluster) (api.Cluster, error)

var (
	clusterFactory = make(map[config.ClusterType]ClusterCreator)
)

func registCluster(t config.ClusterType, creator ClusterCreator) {
	clusterFactory[t] = creator
}

func NewCluster(c config.Cluster) (api.Cluster, error) {
	if creator, ok := clusterFactory[c.GetType()]; ok {
		return creator(c)
	}
	return nil, fmt.Errorf("not support cluster type: %v", c.GetType())
}

type SimpleCluster struct {
	snapShot atomic.Value
	info     config.Cluster
	update   time.Time
	qr       *qos.QosRouterImpl
}

func (c *SimpleCluster) Snapshot() api.ClusterSnapshot {
	ss := c.snapShot.Load()
	if css, ok := ss.(*ClusterSnapShotImpl); ok {
		return css
	}
	return nil
}

func (c *SimpleCluster) GetQosRouter() api.QosRouter {
	return c.qr
}

func (c *SimpleCluster) UpdateTime() time.Time {
	return c.update
}

func (c *SimpleCluster) Config() interface{} {
	return c.info
}

func (c *SimpleCluster) UpdateEndpoints(edps config.LbEndpointSet) {
	// 此处不对cluster info设置endpoints了，而是单独保存了一份，
	// 一来是防止竞争，二来认为cluterinfo是静态的数据
	snapShot := &ClusterSnapShotImpl{
		clusterInfo: c.info,
		// lb:          lb.NewLoadBalancer(c.info.GetLbPolicy(), edps),
		endPoints: edps}
	c.snapShot.Store(snapShot)
	c.qr.UpdateEndpoints(edps)
}

func (c *SimpleCluster) getConfigEndpoints() (config.LbEndpointSet, error) {
	var edpSet config.LbEndpointSet
	for _, edp := range c.info.GetLoadAssignment().GetEndpoints() {
		for _, lbe := range edp.GetLbEndpoint() {
			if lbe.GetLoadBalancingWeight() == 0 {
				lbe.SetLoadBalancingWeight(config.Default_LoadBalance_Weight)
			}
			edpSet = append(edpSet, lbe)
		}
	}
	return edpSet, nil
}

func newSimpleCluster(cluster config.Cluster) *SimpleCluster {

	clusterType := cluster.GetType()
	lbType := cluster.GetLbPolicy()
	if clusterType == config.Cluster_ORIGINAL_DST {
		lbType = config.LB_Original_Dst
	} else if clusterType == config.Cluster_Logical_DNS {
		lbType = config.LB_Logical_DNS
	}
	log.Debug("cluster:%s type:%d lb:%d", cluster.GetName(), clusterType, lbType)

	cluster.SetLbPolicy(lbType)
	return &SimpleCluster{
		info:   cluster,
		update: time.Now(),
		qr:     qos.NewQosRouter(nil, lbType, nil, nil)}
}

func (c *SimpleCluster) Close() {}

type ClusterSnapShotImpl struct {
	clusterInfo config.Cluster
	// lb          api.LoadBalancer
	endPoints config.LbEndpointSet
}

func (cs *ClusterSnapShotImpl) ClusterInfo() config.Cluster {
	return cs.clusterInfo
}

func (cs *ClusterSnapShotImpl) EndpointSet() config.LbEndpointSet {
	return cs.endPoints
}

// func (cs *ClusterSnapShotImpl) LoadBalancer() api.LoadBalancer {
// 	return cs.lb
// }
