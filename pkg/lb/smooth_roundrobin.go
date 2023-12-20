package lb

import (
	"sync"

	"github.com/wereliang/sota-mesh/pkg/api"
	"github.com/wereliang/sota-mesh/pkg/config"
)

type rbItem struct {
	endpoint config.LbEndpoint
	step     int32
}

type SmoothRoundRobin struct {
	sync.RWMutex
	total int32
	items []*rbItem
}

func (lb *SmoothRoundRobin) Select(api.LoadBalancerContext) config.LbEndpoint {
	lb.RLock()
	defer lb.RUnlock()

	maxIndex := 0
	for i := 0; i < len(lb.items); i++ {
		item := lb.items[i]
		item.step += int32(item.endpoint.GetLoadBalancingWeight())
		if item.step > lb.items[maxIndex].step {
			maxIndex = i
		}
		// log.Trace("host:%s weight:%d", item.host.Address().String(), item.step)
	}
	lb.items[maxIndex].step -= lb.total
	return lb.items[maxIndex].endpoint
}

func init() {
	registLoadBalancer(config.Round_Robin, func(edps config.LbEndpointSet) api.LoadBalancer {
		rb := &SmoothRoundRobin{}
		for _, e := range edps {
			rb.total += int32(e.GetLoadBalancingWeight())
			rb.items = append(rb.items, &rbItem{endpoint: e})
		}
		return rb
	})
}
