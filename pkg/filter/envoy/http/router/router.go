package httprouter

import (
	"net"
	"time"

	envoy_extensions_filters_http_router_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/router/v3"

	"github.com/wereliang/sota-mesh/pkg/api"
	"github.com/wereliang/sota-mesh/pkg/config"
	"github.com/wereliang/sota-mesh/pkg/dispatch/http"
	"github.com/wereliang/sota-mesh/pkg/filter"
	"github.com/wereliang/sota-mesh/pkg/log"
	"github.com/wereliang/sota-mesh/pkg/qos"
)

func init() {
	filter.HTTPFilterFactory.Regist(new(RouterFactory))
}

type Router struct {
	cb      api.DecoderFilterCallbacks
	context api.FactoryContext
}

func (r *Router) SetDecoderFilterCallbacks(cb api.DecoderFilterCallbacks) {
	r.cb = cb
}

func (r *Router) Decode(ctx api.StreamContext) api.FilterStatus {
	route := r.cb.Route()
	entry := route.Match(ctx.Request().Header())
	if entry == nil {
		log.Error("route match fail")
		return api.Stop
	}

	log.Debug("[Cluster: %s]", entry.GetClusterName())

	cluster := r.context.ClusterManager().GetCluster(entry.GetClusterName())
	if cluster == nil {
		log.Error("not found cluster:%s", entry.GetClusterName())
		return api.Stop
	}

	_, err := qos.QosCall(cluster, r.cb, func(qosRes api.QosResult) (interface{}, error) {
		ctx.Request().SetHost(qosRes.Endpoint.GetEndpoint().GetAddress().String())
		return nil, http.CallTimeout(ctx, r.getSourceAddr(cluster.Snapshot().ClusterInfo()), time.Second*10)
	})
	if err != nil {
		log.Error("http call error: %s", err)
		return api.Stop
	}

	return api.Continue
}

func (r *Router) getSourceAddr(cluster config.Cluster) net.Addr {
	if bind := cluster.GetUpstreamBindConfig(); bind != nil {
		if addr := bind.GetSourceAddress(); addr != nil {
			return &net.TCPAddr{IP: net.ParseIP(addr.Address), Port: int(addr.PortValue)}
		}
	}
	return nil
}

func (r *Router) Encode(ctx api.StreamContext) api.FilterStatus {
	return api.Continue
}

type RouterFactory struct {
}

func (f *RouterFactory) Name() string {
	return filter.HTTP_Router
}

func (f *RouterFactory) CreateEmptyConfigProto() interface{} {
	return &envoy_extensions_filters_http_router_v3.Router{}
}

func (f *RouterFactory) CreateFilterFactory(pb interface{}, context api.FactoryContext) api.HTTPFilterCreator {
	return func(cb api.HTTPFilterManager) {
		router := &Router{context: context}
		cb.AddDecodeFilter(router)
		cb.AddEncodeFilter(router)
	}
}
