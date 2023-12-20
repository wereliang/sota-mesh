package simpleline

import (
	"bytes"
	"fmt"

	"github.com/wereliang/sota-mesh/pkg/api"
	"github.com/wereliang/sota-mesh/pkg/dispatch"
	"github.com/wereliang/sota-mesh/pkg/dispatch/simple"
	"github.com/wereliang/sota-mesh/pkg/filter"
	"github.com/wereliang/sota-mesh/pkg/log"
)

func init() {
	filter.NetworkFilterFactory.Regist(new(SimpleLineFactory))
}

type SimpleFilter struct {
	dispatcher dispatch.Dispatcher
}

func newSimpleFilter(pb interface{}, cb api.ConnectionCallbacks) api.ReadFilter {
	if config, ok := pb.(*SimpleConfig); !ok {
		panic("simple config error")
	} else {
		log.Debug("simple config: %#v\n", config)
	}
	return &SimpleFilter{dispatcher: simple.NewSimpeDispatcher(cb)}
}

func (z *SimpleFilter) OnData(buf *bytes.Buffer) api.FilterStatus {
	if err := z.dispatcher.Dispatch(buf); err != nil {
		log.Error("%s", err)
		return api.Stop
	}
	return api.Continue
}

// OnNewConnection is called on new connection is created
func (z *SimpleFilter) OnNewConnection() api.FilterStatus {
	log.Debug("call simple on connection")
	return api.Continue
}

type SimpleLineFactory struct {
}

func (f *SimpleLineFactory) Name() string {
	return "sota.filters.network.simple_line"
}

func (f *SimpleLineFactory) CreateEmptyConfigProto() interface{} {
	return &SimpleConfig{}
}

func (f *SimpleLineFactory) CreateFilterFactory(
	pb interface{}, context api.FactoryContext) api.NetworkFilterCreator {

	return func(fm api.FilterManager, cb api.ConnectionCallbacks) error {
		sfilter := newSimpleFilter(pb, cb)
		if sfilter == nil {
			return fmt.Errorf("create simple filter fail")
		}
		fm.AddReadFilter(sfilter)
		return nil
	}
}
