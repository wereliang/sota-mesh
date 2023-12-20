package network

import (
	"net"
	"sync/atomic"

	"github.com/wereliang/sota-mesh/pkg/api"
	"github.com/wereliang/sota-mesh/pkg/log"
)

func NewListener(addr net.Addr) api.Listener {
	l := &listener{addr: addr}
	return l
}

type listener struct {
	addr net.Addr
	cb   atomic.Value
}

func (nl *listener) Listen() error {
	log.Debug("network listen. %s %s", nl.addr.Network(), nl.addr.String())
	l, err := net.Listen(nl.addr.Network(), nl.addr.String())
	if err != nil {
		return err
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}

		go func() {
			cb := nl.cb.Load().(api.ListenerCallback)
			cb.OnAccept(newConnSize(conn, 4096))
		}()
	}
}

func (nl *listener) SetCallback(cb api.ListenerCallback) {
	nl.cb.Store(cb)
}

func (nl *listener) Addr() net.Addr {
	return nl.addr
}
