package http_inspector

import (
	http_inspector_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/listener/http_inspector/v3"
	"github.com/wereliang/sota-mesh/pkg/api"
	"github.com/wereliang/sota-mesh/pkg/filter"
	"github.com/wereliang/sota-mesh/pkg/log"
)

func init() {
	filter.ListenerFilterFactory.Regist(new(HttpInspectorFactory))
}

const (
	HTTP11 = "http/1.1"
)

var (
	minMethodLengh = len("GET")
	maxMethodLengh = len("CONNECT")
	httpMethod     = map[string]struct{}{
		"OPTIONS": {},
		"GET":     {},
		"HEAD":    {},
		"POST":    {},
		"PUT":     {},
		"DELETE":  {},
		"TRACE":   {},
		"CONNECT": {},
	}
)

type HttpInspectorFilter struct {
}

func (f *HttpInspectorFilter) OnAccept(cb api.ListenerFilterCallbacks) api.FilterStatus {
	c := cb.Connection()
	data, err := c.Peek(maxMethodLengh)
	if err != nil {
		log.Error("peer error. %s", err)
		return api.Continue
	}

	size := len(data)
	if size < minMethodLengh {
		log.Error("peer error. %s", err)
		return api.Continue
	}

	if size > maxMethodLengh {
		size = maxMethodLengh
	}

	// check http1, 这里不是很严谨
	for i := minMethodLengh; i <= size; i++ {
		if _, ok := httpMethod[string(data[:i])]; ok {
			cb.Connection().Context().SetApplicationProtocol(HTTP11)
			log.Debug("check http1.1 by method:%s", string(data[:i]))
			return api.Continue
		}
	}

	log.Error("check http fail")
	return api.Continue
}

type HttpInspectorFactory struct {
}

func (f *HttpInspectorFactory) Name() string {
	return filter.Listener_HttpInspector
}

func (f *HttpInspectorFactory) CreateEmptyConfigProto() interface{} {
	return &http_inspector_v3.HttpInspector{}
}

func (f *HttpInspectorFactory) CreateFilterFactory(pb interface{}, context api.FactoryContext) api.ListenerFilterCreator {
	return func(cb api.ListenerFilterManager) {
		cb.AddAcceptFilter(&HttpInspectorFilter{})
	}
}
