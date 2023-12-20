package listener

import (
	"fmt"
	"net"
	"time"

	"github.com/wereliang/sota-mesh/pkg/api"
	"github.com/wereliang/sota-mesh/pkg/config"
	"github.com/wereliang/sota-mesh/pkg/filter"
	"github.com/wereliang/sota-mesh/pkg/log"
	"github.com/wereliang/sota-mesh/pkg/network"
)

func NewActiveListener(
	ltype config.ListenerType, config config.Listener,
	context api.FactoryContext, netl api.Listener) (api.ActiveListener, error) {

	if netl == nil {
		netl = network.NewListener(config.GetAddress())
	}

	fcm := NewFilterChainManager(config.GetFilterChains(), config.GetDefaultFilterChain())

	al := &ActiveListenerImpl{
		config:             config,
		context:            context,
		ltype:              ltype,
		listener:           netl,
		filterChainManager: fcm,
		useOriginalDst:     config.GetUseOriginalDst(),
		bindToPort:         config.GetBindToPort(),
		update:             time.Now(),
	}

	// l.SetCallback(al)
	return al, nil
}

type ActiveListenerImpl struct {
	config             config.Listener
	ltype              config.ListenerType
	context            api.FactoryContext
	filters            []api.ListenerFilter
	listener           api.Listener
	filterChainManager api.FilterChainManager
	useOriginalDst     bool
	bindToPort         bool
	update             time.Time
}

func (l *ActiveListenerImpl) Start() error {
	return l.listener.Listen()
}

func (l *ActiveListenerImpl) UpdateTime() time.Time {
	return l.update
}

func (l *ActiveListenerImpl) Config() interface{} {
	return l.config
}

func (l *ActiveListenerImpl) Type() config.ListenerType {
	return l.ltype
}

func (l *ActiveListenerImpl) Listener() api.Listener {
	return l.listener
}

func (l *ActiveListenerImpl) GetUseOriginalDst() bool {
	return l.useOriginalDst
}

func (l *ActiveListenerImpl) GetBindToPort() bool {
	return l.bindToPort
}

func (l *ActiveListenerImpl) AddAcceptFilter(f api.ListenerFilter) {
	l.filters = append(l.filters, f)
}

func (l *ActiveListenerImpl) OnAccept(conn api.Connection) {
	log.Debug("[Listener: %s]", l.config.GetName())

	lcb := &listenerCallbacks{conn}

	if !l.onListenerFilter(lcb) {
		conn.Close()
		return
	}

	if l.GetUseOriginalDst() {
		ip, port := conn.Context().GetDestinationIP(), conn.Context().GetDestinationPort()
		rdl := l.getRedirectListener(ip, port)
		//  If there is no listener associated with the original destination address,
		//  the connection is handled by the listener that receives it
		if rdl != nil {
			log.Debug("redirect listener to: %s", rdl.Listener().Addr().String())
			cb := rdl.(api.ListenerCallback)
			cb.OnAccept(conn)
			return
		}
		log.Debug("get redirect listener fail: %v %d", ip, port)
	}

	ac := NewActiveConnection(conn)
	filters := l.matchFilters(conn.Context())
	if filters == nil {
		log.Error("Match filter chain fail")
		conn.Close()
		return
	}

	for _, f := range filters {
		// TODO: 此处有点耗性能，提前添加
		factory, pb := filter.GetNetworkFactory(f.GetTypedConfig(), f.GetName())
		if factory == nil {
			if filter.IsWellknowName(f.GetName()) {
				panic(fmt.Errorf("not found network factory:%s", f.GetName()))
			} else {
				log.Error("not support network filter: %s", f.GetName())
				continue
			}
		}
		if err := factory.CreateFilterFactory(pb, l.context)(ac, ac); err != nil {
			log.Error("create network filter fail: %s", err)
			conn.Close()
			return
		}
	}

	ac.OnLoop()
}

func (l *ActiveListenerImpl) onListenerFilter(cb api.ListenerFilterCallbacks) bool {
	for _, f := range l.filters {
		if f.OnAccept(cb) == api.Stop {
			return false
		}
	}
	return true
}

func (l *ActiveListenerImpl) matchFilters(cs api.ConnectionContext) []config.Filter {
	filterChain := l.filterChainManager.FindFilterChains(cs)
	if filterChain != nil {
		log.Debug("[FilterChain: %s]", filterChain.GetName())
		return filterChain.GetFilters()
	}
	return nil
}

func (l *ActiveListenerImpl) getRedirectListener(ip net.IP, port uint32) api.ActiveListener {
	addr := &net.TCPAddr{IP: ip, Port: int(port)}
	return l.context.ListenerManager().FindListenerByAddress(addr)
}

type listenerCallbacks struct {
	conn api.Connection
}

func (cb *listenerCallbacks) Connection() api.Connection {
	return cb.conn
}
