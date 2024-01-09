package router

import (
	"github.com/wereliang/sota-mesh/pkg/api"
	"github.com/wereliang/sota-mesh/pkg/filter"
	"github.com/wereliang/sota-mesh/pkg/log"
)

func init() {
	filter.HTTPFilterFactory.Regist(new(RouterFactory))
}

type SotaRouterFilter struct {
}

func (r *SotaRouterFilter) SetDecoderFilterCallbacks(cb api.DecoderFilterCallbacks) {
}

func (r *SotaRouterFilter) Decode(ctx api.StreamContext) api.FilterStatus {
	log.Debug("call sota router filter")
	return api.Continue
}

func (r *SotaRouterFilter) Encode(ctx api.StreamContext) api.FilterStatus {
	return api.Continue
}

type RouterFactory struct {
}

func (f *RouterFactory) Name() string {
	return "sota.filters.http.router"
}

func (f *RouterFactory) CreateEmptyConfigProto() interface{} {
	return nil
}

func (f *RouterFactory) CreateFilterFactory(pb interface{}, context api.FactoryContext) api.HTTPFilterCreator {
	return func(cb api.HTTPFilterManager) {
		router := &SotaRouterFilter{}
		cb.AddDecodeFilter(router)
		cb.AddEncodeFilter(router)
	}
}
