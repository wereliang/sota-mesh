package qos

import (
	"github.com/wereliang/sota-mesh/pkg/api"
	"github.com/wereliang/sota-mesh/pkg/log"
)

type QosFunc func(api.QosResult) (interface{}, error)

func QosCall(cluster api.Cluster, lbCtx api.LoadBalancerContext, fn QosFunc) (interface{}, error) {
	qr := cluster.GetQosRouter().(*QosRouterImpl)
	cncy := qr.GetCircuitBreaker("").Concurrency()
	if !cncy.Access() {
		log.Error("ErrTooManyConcurrency")
		return nil, ErrTooManyConcurrency
	}
	cncy.In()
	defer cncy.Out()

	res, err := qr.GetRoute(lbCtx)
	if err != nil {
		return nil, err
	}
	cb := qr.GetCircuitBreaker(hashEndpoint(res.Endpoint))
	return cb.Execute(func() (interface{}, error) {
		return fn(res)
	})
}
