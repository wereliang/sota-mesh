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

	if len(lb.items) == 0 {
		return nil
	}

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

func (lb *SmoothRoundRobin) Refresh(edps []config.LbEndpoint) {
	lb.Lock()
	defer lb.Unlock()
	lb.total = 0
	lb.items = nil
	for _, e := range edps {
		lb.total += int32(e.GetLoadBalancingWeight())
		lb.items = append(lb.items, &rbItem{endpoint: e})
	}
}

func init() {
	registLoadBalancer(config.Round_Robin, func(edps config.LbEndpointSet) api.LoadBalancer {
		rb := &SmoothRoundRobin{}
		rb.Refresh(edps)
		return rb
	})
}
