package zmq

import (
	"bytes"
	"fmt"

	"github.com/wereliang/sota-mesh/pkg/api"
	"github.com/wereliang/sota-mesh/pkg/dispatch/zmtp"
	"github.com/wereliang/sota-mesh/pkg/filter"
	"github.com/wereliang/sota-mesh/pkg/log"
)

func init() {
	filter.NetworkFilterFactory.Regist(new(ZmqFactory))
}

type ZmqFilter struct {
	dispatcher *zmtp.ZmqDispatcher
}

func newZmqFilter(pb interface{}, cb api.ConnectionCallbacks) api.ReadFilter {
	return &ZmqFilter{dispatcher: zmtp.NewZmqDispatcher(cb).(*zmtp.ZmqDispatcher)}
}

func (z *ZmqFilter) OnData(buf *bytes.Buffer) api.FilterStatus {
	if err := z.dispatcher.Dispatch(buf); err != nil {
		log.Error("%s", err)
		return api.Stop
	}
	return api.Continue
}

// zmq有greet的阶段，不能在创建dispatch的时候hold住执行，可以在此处操作
func (z *ZmqFilter) OnNewConnection() api.FilterStatus {
	log.Debug("call zmq on connection")
	go func() { z.dispatcher.Serve() }()
	return api.Continue
}

type ZmqFactory struct {
}

func (f *ZmqFactory) Name() string {
	return "sota.filters.network.zmq"
}

func (f *ZmqFactory) CreateEmptyConfigProto() interface{} {
	return nil
}

func (f *ZmqFactory) CreateFilterFactory(
	pb interface{}, context api.FactoryContext) api.NetworkFilterCreator {

	return func(fm api.FilterManager, cb api.ConnectionCallbacks) error {
		zfilter := newZmqFilter(pb, cb)
		if zfilter == nil {
			return fmt.Errorf("create zmq filter fail")
		}
		fm.AddReadFilter(zfilter)
		return nil
	}
}
