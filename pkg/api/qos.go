package api

import (
	"github.com/wereliang/sota-mesh/pkg/config"
)

type QosResult struct {
	Endpoint config.LbEndpoint
	IsDetect bool
}

type QosRouter interface {
	GetRoute(LoadBalancerContext) (QosResult, error)
	// Report(QosResult, error, time.Duration)
}
