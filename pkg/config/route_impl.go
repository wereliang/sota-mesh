package config

type RouteConfigurationImpl struct {
	Name        string             `json:"name"`
	VirtualHost []*VirtualHostImpl `json:"virtual_hosts"`
}

func (r *RouteConfigurationImpl) GetName() string {
	return r.Name
}

func (r *RouteConfigurationImpl) GetVirtualHosts() []VirtualHost {
	var vhs []VirtualHost
	for _, vh := range r.VirtualHost {
		vhs = append(vhs, vh)
	}
	return vhs
}

type VirtualHostImpl struct {
	Name    string       `json:"name"`
	Domains []string     `json:"domains"`
	Routes  []*RouteImpl `json:"routes"`
}

func (v *VirtualHostImpl) GetName() string {
	return v.Name
}

func (v *VirtualHostImpl) GetDomains() []string {
	return v.Domains
}

func (v *VirtualHostImpl) GetRoutes() []Route {
	var routes []Route
	for _, r := range v.Routes {
		routes = append(routes, r)
	}
	return routes
}

type RouteImpl struct {
	Match *RouteMatchImpl  `json:"match"`
	Route *RouteActionImpl `json:"route"`
}

func (r *RouteImpl) GetMatch() RouteMatch {
	return r.Match
}

func (r *RouteImpl) GetRoute() RouteAction {
	return r.Route
}

type RouteMatchImpl struct {
	Prefix string `json:"prefix"`
	Path   string `json:"path"`
}

func (r *RouteMatchImpl) GetPath() string {
	return r.Path
}

func (r *RouteMatchImpl) GetPrefix() string {
	return r.Prefix
}

type RouteActionImpl struct {
	Cluster string `json:"cluster"`
}

func (r *RouteActionImpl) GetCluster() string {
	return r.Cluster
}

type RouteEntryImpl struct {
	ClusterName string
}

func (re *RouteEntryImpl) GetClusterName() string {
	return re.ClusterName
}
