package route

import (
	"fmt"

	"github.com/fanyang01/radix"
	"github.com/wereliang/sota-mesh/pkg/api"
	"github.com/wereliang/sota-mesh/pkg/config"
	"github.com/wereliang/sota-mesh/pkg/log"
)

// https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/route/v3/route_components.proto
// Domain search order:
// Exact domain names: www.foo.com.
// Suffix domain wildcards: *.foo.com or *-bar.foo.com.
// Prefix domain wildcards: foo.* or foo-*.
// Special wildcard * matching any domain.

type RouteConfigMatcherImpl struct {
	config  config.RouteConfiguration
	domains *radix.PatternTrie
}

func NewRouterMatcher(c config.RouteConfiguration) api.RouteConfigMatcher {
	rc := &RouteConfigMatcherImpl{config: c, domains: radix.NewPatternTrie()}
	rc.build()
	return rc
}

func (r *RouteConfigMatcherImpl) build() {

	for _, vhConfig := range r.config.GetVirtualHosts() {
		match := newPathTrie(vhConfig.GetName())
		for _, route := range vhConfig.GetRoutes() {
			cluster := route.GetRoute().GetCluster()
			if cluster == "" {
				panic("invalid route action(just support cluster)")
			}

			if route.GetMatch() == nil {
				log.Error("virtual_hosts:%s config empty match", vhConfig.GetName())
				continue
			}

			if path := route.GetMatch().GetPath(); path != "" {
				match.Path(path, cluster)
			} else if prefix := route.GetMatch().GetPrefix(); prefix != "" {
				match.Prefix(prefix, cluster)
			} else {
				panic(fmt.Sprintf("invalid match path:%#v", route))
			}
		}

		for _, host := range vhConfig.GetDomains() {
			r.domains.Add(host, match)
		}
	}
}

func (r *RouteConfigMatcherImpl) Match(header api.RequestHeader) config.RouteEntry {
	v, ok := r.domains.Lookup(string(header.Host()))
	if !ok {
		return nil
	}
	if re := v.(*PathTrie).Lookup(string(header.Path())); re != nil {
		return &config.RouteEntryImpl{ClusterName: re.(string)}
	}
	return nil
}

type PathTrie struct {
	Name string
	*radix.PatternTrie
}

func newPathTrie(name string) *PathTrie {
	return &PathTrie{Name: name, PatternTrie: &radix.PatternTrie{}}
}

func (t *PathTrie) Path(s string, v interface{}) {
	t.PatternTrie.Add(s, v)
}

func (t *PathTrie) Prefix(s string, v interface{}) {
	t.PatternTrie.Add(s+"*", v)
}

func (t *PathTrie) Lookup(s string) interface{} {
	v, _ := t.PatternTrie.Lookup(s)
	return v
}
