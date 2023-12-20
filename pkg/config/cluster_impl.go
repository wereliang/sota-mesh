package config

type ClusterImpl struct {
	Name               string                     `json:"name"`
	LoadAssignment     *ClusterLoadAssignmentImpl `json:"load_assignment"`
	LbPolicy           LoadBalancerType           `json:"lb_policy"`
	UpstreamBindConfig *BindConfigImpl            `json:"upstream_bind_config"`
	ClusterType        ClusterType
}

func (c *ClusterImpl) GetName() string {
	return c.Name
}

func (c *ClusterImpl) GetType() ClusterType {
	return c.ClusterType
}

func (c *ClusterImpl) GetLbPolicy() LoadBalancerType {
	return c.LbPolicy
}

func (c *ClusterImpl) SetLbPolicy(lbt LoadBalancerType) {
	c.LbPolicy = lbt
}

func (c *ClusterImpl) GetLoadAssignment() ClusterLoadAssignment {
	return c.LoadAssignment
}

func (c *ClusterImpl) GetUpstreamBindConfig() BindConfig {
	if c.UpstreamBindConfig != nil {
		return c.UpstreamBindConfig
	}
	return nil
}

type ClusterLoadAssignmentImpl struct {
	ClusterName string                     `json:"cluster_name"`
	Endpoints   []*LocalityLbEndpointsImpl `json:"endpoints"`
	endpoints   []LocalityLbEndpoints
}

func (c *ClusterLoadAssignmentImpl) GetClusterName() string {
	return c.ClusterName
}

func (c *ClusterLoadAssignmentImpl) GetEndpoints() []LocalityLbEndpoints {
	if c.endpoints == nil {
		for _, edp := range c.Endpoints {
			c.endpoints = append(c.endpoints, edp)
		}
	}
	return c.endpoints

	// var endpoints []LocalityLbEndpoints
	// for _, edp := range c.Endpoints {
	// 	endpoints = append(endpoints, edp)
	// }
	// return endpoints
}
