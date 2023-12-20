package listener

import (
	"fmt"
	"net"
	"sync"

	"github.com/wereliang/sota-mesh/pkg/api"
	"github.com/wereliang/sota-mesh/pkg/config"
	"github.com/wereliang/sota-mesh/pkg/filter"
	"github.com/wereliang/sota-mesh/pkg/log"
)

type ListenerManagerImpl struct {
	listenerMap sync.Map
	context     api.FactoryContext
}

func NewListenerManager(ctx api.FactoryContext) (api.ListenerManager, error) {
	return &ListenerManagerImpl{context: ctx}, nil
}

func (lm *ListenerManagerImpl) AddOrUpdateListener(ltype config.ListenerType, cfg config.Listener) error {
	var (
		actl   api.ActiveListener
		netl   api.Listener
		err    error
		update bool = false
	)

	if any, ok := lm.listenerMap.Load(cfg.GetName()); ok {
		netl = any.(api.ActiveListener).Listener()
		// addr must be same
		if netl.Addr().String() != cfg.GetAddress().String() {
			return fmt.Errorf("addr not same. (%s) != (%s)",
				netl.Addr().String(), cfg.GetAddress().String())
		}
		update = true
	}

	if actl, err = NewActiveListener(ltype, cfg, lm.context, netl); err != nil {
		return err
	}
	if err = lm.addListenerFilter(cfg.GetListenerFilters(), actl); err != nil {
		return err
	}

	actl.Listener().SetCallback(actl.(api.ListenerCallback))
	lm.listenerMap.Store(cfg.GetName(), actl)

	// TODO: stop and destory
	if !update && actl.GetBindToPort() {
		go func() {
			if err := actl.Start(); err != nil {
				panic(err)
			}
		}()
	}

	if len(cfg.GetListenerFilters()) > 0 {
		lfilter := cfg.GetListenerFilters()[0]
		log.Debug("listener filter: %#v type:%#v", lfilter, lfilter.GetTypedConfig())
	}

	return nil
}

// FindListenerByAddress find listerner by address
func (lm *ListenerManagerImpl) FindListenerByAddress(net.Addr) api.ActiveListener {
	return nil
}

// FindListenerByName find listener by name
func (lm *ListenerManagerImpl) FindListenerByName(string) api.ActiveListener {
	return nil
}

// DeleteListener delete exist listener
func (lm *ListenerManagerImpl) DeleteListener(string) error {
	return nil
}

func (lm *ListenerManagerImpl) addListenerFilter(
	filters []config.Filter, actl api.ActiveListener) error {

	for _, f := range filters {
		factory, pb := filter.GetListenerFactory(f.GetTypedConfig(), f.GetName())
		if factory == nil {
			log.Error("not support listener filter (%s) now", f.GetName())
			// if filter.IsWellknowName(f.GetName()) {
			// 	panic(fmt.Errorf("not found listener factory:%s", f.Name))
			// } else {
			// 	log.Error("not support listener filter (%s) now", f.Name)
			// 	continue
			// }
			continue
		}
		factory.CreateFilterFactory(pb, lm.context)(actl)
	}

	// add original dst filter if use_original_dst flag set
	if actl.GetUseOriginalDst() {
		return lm.buildOriginalDstListenerFilter(actl)
	}
	return nil
}

func (lm *ListenerManagerImpl) buildOriginalDstListenerFilter(actl api.ActiveListener) error {
	factory, pb := filter.GetListenerFactory(nil, filter.Listener_OriginalDst)
	if factory == nil {
		return fmt.Errorf("not found factory:%s", filter.Listener_OriginalDst)
	}
	factory.CreateFilterFactory(pb, lm.context)(actl)
	return nil
}
