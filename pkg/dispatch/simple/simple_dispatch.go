package simple

import (
	"fmt"
	"net/url"

	"github.com/wereliang/sota-mesh/pkg/api"
	"github.com/wereliang/sota-mesh/pkg/dispatch"
	"github.com/wereliang/sota-mesh/pkg/log"
)

type SimpleDispatcher struct {
	*dispatch.DispatcherImpl
}

func NewSimpeDispatcher(conn api.ConnectionCallbacks) dispatch.Dispatcher {
	dispatch, _ := dispatch.NewDispatcher(conn)
	dp := &SimpleDispatcher{DispatcherImpl: dispatch}
	go dp.serve()
	return dp
}

func (s *SimpleDispatcher) serve() {
	defer s.Close()
	for {
		line, _, err := s.Reader.ReadLine()
		if err != nil {
			log.Debug("close by peer")
			return
		}
		log.Debug("receive: %s", string(line))
		vs, err := url.ParseQuery(string(line))
		if err != nil {
			s.Conn.Write([]byte(err.Error()))
			return
		}
		s.Conn.Write([]byte(fmt.Sprintf("receive:%#v\r\n", vs)))
	}
}
