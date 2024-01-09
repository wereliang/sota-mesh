package hcm

import (
	"bytes"
	"fmt"

	envoy_filters_network_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"

	"github.com/wereliang/sota-mesh/pkg/api"
	"github.com/wereliang/sota-mesh/pkg/config"
	"github.com/wereliang/sota-mesh/pkg/config/envoyv3"
	"github.com/wereliang/sota-mesh/pkg/dispatch"
	"github.com/wereliang/sota-mesh/pkg/dispatch/http"
	"github.com/wereliang/sota-mesh/pkg/filter"
	"github.com/wereliang/sota-mesh/pkg/log"
	"github.com/wereliang/sota-mesh/pkg/route"
)

func init() {
	filter.NetworkFilterFactory.Regist(new(HttpConnectionManagerFactory))
}

type HttpConnectionManager struct {
	api.ReadFilter
	dispatcher dispatch.Dispatcher
	config     *envoy_filters_network_v3.HttpConnectionManager
}

func getRouteConfiguration(hcm *envoy_filters_network_v3.HttpConnectionManager,
	rcm api.RouteConfigManager) config.RouteConfiguration {

	if c := hcm.GetRouteConfig(); c != nil {
		return &envoyv3.RouteConfigurationEnvoy{RouteConfiguration: c}
	}
	if hcm.GetRds() != nil {
		if rc := rcm.GetRouteConfig(hcm.GetRds().GetRouteConfigName()); rc != nil {
			return rc
		}
		log.Error("not found rds route config [%s]", hcm.GetRds().GetRouteConfigName())
		return nil
	}
	// Not support other type now
	panic("invalid route config")
}

func newHttpConnectionManager(pb interface{}, cb api.ConnectionCallbacks, context api.FactoryContext) api.ReadFilter {
	cfg := pb.(*envoy_filters_network_v3.HttpConnectionManager)
	log.Debug("HttpConnectionManager config: %#v", cfg)
	hcm := &HttpConnectionManager{config: cfg}

	// TODO: 复用
	rc := getRouteConfiguration(cfg, context.RouteConfigManager())
	if rc == nil {
		return nil
	}
	log.Debug("[RouteConfig: %s]", rc.GetName())

	matcher := route.NewRouterMatcher(rc)
	handler := http.NewHandler(matcher, cb)

	for _, f := range hcm.config.HttpFilters {
		var any config.TypeConfig
		if f.GetTypedConfig() != nil {
			any = &config.AnyTypeConfig{A: f.GetTypedConfig()}
		}
		factory, pb := filter.GetHTTPFactory(any, f.Name)
		if factory == nil {
			if filter.IsWellknowName(f.Name) {
				panic(fmt.Errorf("not found factory:%s", f.Name))
			} else {
				log.Error("not support http filter: %s", f.Name)
				continue
			}
		}
		factory.CreateFilterFactory(pb, context)(handler)
	}
	hcm.dispatcher = http.NewStreamServer(handler, cb)
	return hcm
}

func (f *HttpConnectionManager) OnData(buffer *bytes.Buffer) api.FilterStatus {
	if err := f.dispatcher.Dispatch(buffer); err != nil {
		return api.Stop
	}
	return api.Continue
}

func (f *HttpConnectionManager) OnNewConnection() api.FilterStatus {
	return api.Continue
}

type HttpConnectionManagerFactory struct {
}

func (f *HttpConnectionManagerFactory) Name() string {
	return filter.Network_HttpConnectionManager
}

func (f *HttpConnectionManagerFactory) CreateEmptyConfigProto() interface{} {
	return &envoy_filters_network_v3.HttpConnectionManager{}
}

func (f *HttpConnectionManagerFactory) CreateFilterFactory(
	pb interface{}, context api.FactoryContext) api.NetworkFilterCreator {

	return func(fm api.FilterManager, cb api.ConnectionCallbacks) error {
		hcm := newHttpConnectionManager(pb, cb, context)
		if hcm == nil {
			return fmt.Errorf("create http connection manager fail")
		}
		fm.AddReadFilter(hcm)
		return nil
	}
}
