package http

import (
	"context"
	"io"
	"net"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/wereliang/sota-mesh/pkg/api"
	"github.com/wereliang/sota-mesh/pkg/dispatch"
	"github.com/wereliang/sota-mesh/pkg/log"
)

func NewStreamServer(handler Handler, conn api.ConnectionCallbacks) dispatch.Dispatcher {
	dispatcher, _ := dispatch.NewDispatcher(conn)

	s := &httpStreamServer{
		DispatcherImpl: dispatcher,
		handler:        handler,
	}
	go func() {
		s.serve()
	}()
	return s
}

type httpStreamServer struct {
	*dispatch.DispatcherImpl
	handler Handler
}

func (s *httpStreamServer) serve() {
	for {
		request, response := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
		// blocking read using fasthttp.Request.Read
		err := request.ReadLimitBody(s.Reader, 1024*1024)
		if err != nil {
			if err != io.EOF {
				log.Error("ReadLimitBody error: %s. conn close", err)
				s.Close()
			}
			break
		}

		ctx := NewStreamContext(context.TODO(), request, response)
		err = s.handle(ctx)
		if err != nil {
			log.Error("handle error : %s", err)
		}
		response.WriteTo(s.Conn)
	}
	log.Debug("server close")
}

func (s *httpStreamServer) handle(ctx api.StreamContext) error {
	if err := s.handler.Decode(ctx); err != nil {
		return err
	}
	if err := s.handler.Encode(ctx); err != nil {
		return err
	}
	return nil
}

type StreamClient interface {
	Call(api.StreamContext, time.Duration) error
}

// 此处只区分了127.0.0.6的特殊情况
func NewStreamClient(addr net.Addr) StreamClient {
	var sc *fasthttpStreamClient
	if addr == nil {
		sc = &fasthttpStreamClient{}
	} else {
		sc = &fasthttpStreamClient{
			client: fasthttp.Client{Dial: Dial},
		}
	}
	sc.client.MaxIdleConnDuration = time.Minute
	sc.client.DisableHeaderNamesNormalizing = true
	sc.client.MaxConnsPerHost = 30000
	return sc
}

type fasthttpStreamClient struct {
	client fasthttp.Client
}

func (sc *fasthttpStreamClient) Call(ctx api.StreamContext, t time.Duration) error {
	request := ctx.Request().Raw().(*fasthttp.Request)
	response := ctx.Response().Raw().(*fasthttp.Response)
	request.UseHostHeader = true
	return sc.client.DoTimeout(request, response, t)
}

var dialerWithLAddr = &fasthttp.TCPDialer{
	Concurrency: 1000, LocalAddr: &net.TCPAddr{IP: net.ParseIP("127.0.0.6")}}

func Dial(addr string) (net.Conn, error) {
	return dialerWithLAddr.Dial(addr)
}

var defaultStreamClient = NewStreamClient(nil)
var streamClientWithLAddr = NewStreamClient(&net.TCPAddr{})

func Call(ctx api.StreamContext, laddr net.Addr) error {
	return CallTimeout(ctx, laddr, time.Second*5)
}

func CallTimeout(ctx api.StreamContext, laddr net.Addr, t time.Duration) error {
	if laddr != nil {
		return streamClientWithLAddr.Call(ctx, t)
	}
	return defaultStreamClient.Call(ctx, t)
}
