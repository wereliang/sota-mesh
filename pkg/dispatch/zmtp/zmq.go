package zmtp

import (
	"fmt"

	"github.com/wereliang/sota-mesh/pkg/api"
	"github.com/wereliang/sota-mesh/pkg/dispatch"
	"github.com/wereliang/sota-mesh/pkg/log"
)

type ZmqDispatcher struct {
	*dispatch.DispatcherImpl
	zmqSock *Socket
}

func NewZmqDispatcher(conn api.ConnectionCallbacks) dispatch.Dispatcher {

	dispatch, _ := dispatch.NewDispatcher(conn)
	dp := &ZmqDispatcher{DispatcherImpl: dispatch}

	// sock, err := socketFromConnection(conn.Raw().(*net.TCPConn), &SocketConfig{Type: ROUTER})
	sock, err := socketFromReadWriter(dp, &SocketConfig{Type: ROUTER})
	if err != nil {
		panic(err)
	}

	dp.zmqSock = sock
	return dp
}

func (dp *ZmqDispatcher) Serve() {

	s := dp.zmqSock

	err := s.prepare()
	if err != nil {
		s.Close()
		log.Error("prepare error. %s", err)
		return
	}
	s.startChannels(nil)

	for {
		v := <-dp.zmqSock.Read()
		if dp.zmqSock.Error() != nil {
			dp.zmqSock.Close()
			fmt.Println(dp.zmqSock.Error())
			return
		}

		// REQ第一帧为空分隔符帧
		for _, body := range v {
			if len(body) == 0 {
				continue
			}
			log.Debug("zmq recv:%s", string(body))
		}

		dp.zmqSock.Send(v)
	}
}
