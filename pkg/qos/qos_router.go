package qos

import (
	"container/list"
	"fmt"
	"sync"
	"time"

	"github.com/wereliang/sota-mesh/pkg/api"
	"github.com/wereliang/sota-mesh/pkg/config"
	"github.com/wereliang/sota-mesh/pkg/lb"
	"github.com/wereliang/sota-mesh/pkg/log"
)

type DetectFunc func(string) bool

type QosRouterImpl struct {
	sync.RWMutex
	cb            config.CircuitBreakers
	lb            api.LoadBalancer
	closeC        chan struct{}
	healthNodes   NodePairs  // 健康节点
	overloadNodes NodePairs  // 过载节点
	detectList    *list.List // 探测节点
	breakers      CircuitBreakers
	detectFunc    DetectFunc
}

func NewQosRouter(cb config.CircuitBreakers,
	lbType config.LoadBalancerType,
	edps config.LbEndpointSet,
	detectFunc DetectFunc) *QosRouterImpl {
	qr := &QosRouterImpl{
		cb:            cb,
		lb:            lb.NewLoadBalancer(lbType, edps),
		closeC:        make(chan struct{}, 1),
		healthNodes:   newNodeParis(edps),
		overloadNodes: newNodeParis(nil),
		detectList:    list.New(),
		detectFunc:    detectFunc,
	}
	qr.breakers = CircuitBreakers{
		creator: qr.bcCreator,
		checkor: qr.bcCheckor,
	}
	if detectFunc == nil {
		qr.detectFunc = TCPDetect
	}

	go qr.detect()
	return qr
}

func (qr *QosRouterImpl) bcCreator(id string) *CircuitBreaker {
	// todo: from circuit breaker config
	return NewCircuitBreaker(Settings{
		Name:           id,
		MaxRequests:    1,                // 半开状态的探测次数
		Interval:       10 * time.Second, // Close状态的循环计数周期
		Timeout:        30 * time.Second, // Open状态的周期
		MaxConcurrency: 3,
		OnStateChange:  qr.onStateChange,
		ReadyToTrip: func(counts Counts) bool {
			if counts.ConsecutiveFailures > 5 {
				return true
			}
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
	})
}

func (qr *QosRouterImpl) bcCheckor(id string) bool {
	if !qr.healthNodes.ExistByStr(id) && !qr.overloadNodes.ExistByStr(id) {
		return false
	}
	return true
}

func (qr *QosRouterImpl) GetRoute(ctx api.LoadBalancerContext) (api.QosResult, error) {
	detectNode := qr.getDetectNode()
	if detectNode != nil {
		return api.QosResult{Endpoint: detectNode, IsDetect: true}, nil
	}
	if edp := qr.lb.Select(ctx); edp != nil {
		return api.QosResult{Endpoint: edp, IsDetect: false}, nil
	}
	return api.QosResult{}, fmt.Errorf("not found endpoint")
}

func (qr *QosRouterImpl) GetCircuitBreaker(name string) *CircuitBreaker {
	return qr.breakers.GetBreaker(name)
}

func (qr *QosRouterImpl) getDetectNode() config.LbEndpoint {
	qr.RLock()
	if qr.detectList.Len() == 0 {
		qr.RUnlock()
		return nil
	}
	qr.RUnlock()

	qr.Lock()
	defer qr.Unlock()
	for e := qr.detectList.Front(); e != nil; e = e.Next() {
		edp := qr.detectList.Remove(e).(config.LbEndpoint)
		if qr.overloadNodes.Exist(edp) {
			return edp
		}
	}
	return nil
}

func (qr *QosRouterImpl) UpdateEndpoints(edps config.LbEndpointSet) {
	qr.RLock()
	var newOls, newHeals config.LbEndpointSet
	for _, e := range edps {
		if qr.overloadNodes.Exist(e) {
			newOls = append(newOls, e)
			continue
		}
		newHeals = append(newHeals, e)
	}
	qr.RUnlock()

	qr.Lock()
	qr.overloadNodes.Refresh(newOls)
	qr.healthNodes.Refresh(newHeals)
	// 探测节点会在选择的时候做二次校验，因此这里暂不做处理
	qr.lb.Refresh(newHeals)
	// 触发一次circuitbreak清理
	qr.breakers.Clear()
	qr.Unlock()

}

// 获取未探测成功，即未放在探测节点里头的过载节点
func (qr *QosRouterImpl) getDetectOverloadNode() config.LbEndpointSet {
	qr.RLock()
	defer qr.RUnlock()
	olNodes := qr.overloadNodes.edps
	if qr.detectList.Len() == 0 {
		return olNodes
	}

	var resNodes config.LbEndpointSet
	for _, ol := range olNodes {
		exist := false
		for e := qr.detectList.Front(); e != nil; e = e.Next() {
			if e.Value.(config.LbEndpoint).GetEndpoint().GetAddress() ==
				ol.GetEndpoint().GetAddress() {
				exist = true
				break
			}
		}
		if !exist {
			resNodes = append(resNodes, ol)
		}
	}
	return resNodes
}

func (qr *QosRouterImpl) detect() {
	t := time.NewTicker(time.Second * 10)
	defer t.Stop()
	for {
		select {
		case <-qr.closeC:
			return
		case <-t.C:
			nodes := qr.getDetectOverloadNode()
			if len(nodes) == 0 {
				break
			}

			var wg sync.WaitGroup
			resC := make(chan config.LbEndpoint)
			for _, ol := range nodes {
				// 为到达时间的不做处理
				if !ol.GetExpired().Before(time.Now()) {
					continue
				}

				wg.Add(1)
				go func(o config.LbEndpoint) {
					defer wg.Done()
					if ok := qr.detectEndpoint(o); ok {
						resC <- o
					}
				}(ol)
			}

			go func() {
				wg.Wait()
				close(resC)
			}()

			for r := range resC {
				qr.Lock()
				qr.detectList.PushBack(r)
				qr.Unlock()
			}
		}
	}
}

func (qr *QosRouterImpl) setOverload(str string, t time.Time) {
	qr.Lock()
	defer qr.Unlock()
	if e, b := qr.healthNodes.DelByStr(str); b {
		qr.lb.Refresh(qr.healthNodes.edps)
		e.SetExpired(t)
		qr.overloadNodes.Add(e)
		log.Debug("health endpoinds: %s", qr.healthNodes.String())
	}
}

func (qr *QosRouterImpl) recoverOverload(str string) {
	qr.Lock()
	defer qr.Unlock()
	if e, b := qr.overloadNodes.DelByStr(str); b {
		e.SetExpired(time.Time{})
		qr.healthNodes.Add(e)
		qr.lb.Refresh(qr.healthNodes.edps)
		log.Debug("health endpoinds: %s", qr.healthNodes.String())
	}
}

func (qr *QosRouterImpl) detectEndpoint(e config.LbEndpoint) bool {
	return qr.detectFunc(e.GetEndpoint().GetAddress().String())
}

func (qr *QosRouterImpl) Close() {
	qr.closeC <- struct{}{}
}

func (qr *QosRouterImpl) onStateChange(name string, from State, to State, expiry time.Time) {
	log.Debug("onStateChange:%s from:%s to %s", name, from, to)
	if name == "" {
		return
	}
	if to == StateOpen {
		qr.setOverload(name, expiry)
	} else if to == StateClosed {
		qr.recoverOverload(name)
	}
}
