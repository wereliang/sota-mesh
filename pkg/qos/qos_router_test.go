package qos

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wereliang/sota-mesh/pkg/config"
)

func newEndpoint(ip string, port uint32) config.LbEndpoint {
	return &config.LbEndpointImpl{
		Endpoint: &config.EndpointImpl{
			Address: &config.Address{
				SocketAddress: config.SocketAddress{
					Protocol:  "tcp",
					Address:   ip,
					PortValue: port,
				},
			},
		},
		LoadBalancingWeight: 100,
	}
}

func TestUpdateEndpoints(t *testing.T) {
	var edps config.LbEndpointSet
	var (
		e1 = newEndpoint("127.0.0.1", 5566)
		e2 = newEndpoint("127.0.0.1", 7788)
		e3 = newEndpoint("127.0.0.1", 10000)
	)
	edps = append(edps, e1)
	edps = append(edps, e2)
	edps = append(edps, e3)
	qr := NewQosRouter(nil, config.Round_Robin, edps, func(string) bool { return true })

	var b1, b2, b3 bool

	routeFunc := func() {
		b1, b2, b3 = false, false, false
		for i := 0; i < 10; i++ {
			res, err := qr.GetRoute(nil)
			assert.Nil(t, err)
			addr := hashEndpoint(res.Endpoint)
			if addr == hashEndpoint(e1) {
				b1 = true
			} else if addr == hashEndpoint(e2) {
				b2 = true
			} else if addr == hashEndpoint(e3) {
				b3 = true
			}
		}
	}
	routeFunc()

	assert.Equal(t, true, b1)
	assert.Equal(t, true, b2)
	assert.Equal(t, true, b3)

	qr.onStateChange(hashEndpoint(e1), StateClosed, StateOpen, time.Now().Add(time.Second*10))
	qr.onStateChange(hashEndpoint(e2), StateClosed, StateOpen, time.Now().Add(time.Second*10))

	routeFunc()

	assert.Equal(t, false, b1)
	assert.Equal(t, false, b2)
	assert.Equal(t, true, b3)

	// wait for recovery.
	time.Sleep(time.Second * 20)
	routeFunc()
	assert.Equal(t, true, b1)
	assert.Equal(t, true, b2)
	assert.Equal(t, true, b3)
}

func TestQosCall(t *testing.T) {
	var edps config.LbEndpointSet
	var (
		e1 = newEndpoint("127.0.0.1", 5566)
		e2 = newEndpoint("127.0.0.1", 7788)
		e3 = newEndpoint("127.0.0.1", 10000)
	)
	edps = append(edps, e1)
	edps = append(edps, e2)
	edps = append(edps, e3)
	qr := NewQosRouter(nil, config.Round_Robin, edps, func(string) bool { return true })

	// var b1, b2, b3 bool

	// routeFunc := func() {
	// 	b1, b2, b3 = false, false, false
	// 	for i := 0; i < 10; i++ {
	// 		res, err := qr.GetRoute(nil)
	// 		assert.Nil(t, err)
	// 		addr := hashEndpoint(res.Endpoint)
	// 		if addr == hashEndpoint(e1) {
	// 			b1 = true
	// 		} else if addr == hashEndpoint(e2) {
	// 			b2 = true
	// 		} else if addr == hashEndpoint(e3) {
	// 			b3 = true
	// 		}
	// 	}
	// }
	// routeFunc()
	// assert.Equal(t, true, b1)
	// assert.Equal(t, true, b2)
	// assert.Equal(t, true, b3)

	accfn := func() bool {
		cncy := qr.GetCircuitBreaker("").Concurrency()
		if !cncy.Access() {
			return false
		}
		cncy.In()
		defer cncy.Out()
		time.Sleep(time.Second * 3)
		return true
	}

	for i := 0; i < 4; i++ {
		go func() {
			assert.Equal(t, true, accfn())
		}()
	}
	time.Sleep(time.Second)
	assert.Equal(t, false, accfn())

	// test update to clear cb
	qr.GetCircuitBreaker(hashEndpoint(e1))
	qr.GetCircuitBreaker(hashEndpoint(e2))
	qr.GetCircuitBreaker(hashEndpoint(e3))

	edps = edps[:2]
	fmt.Println("edps:", len(edps))
	qr.UpdateEndpoints(edps)

}
