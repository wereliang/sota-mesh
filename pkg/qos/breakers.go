package qos

import (
	"sync"

	"github.com/wereliang/sota-mesh/pkg/log"
)

// type qosRes struct {
// 	res api.QosResult
// 	err error
// 	ts  time.Duration
// }

type CbCreator func(name string) *CircuitBreaker
type CbCheckor func(name string) bool

type CircuitBreakers struct {
	breakers sync.Map
	creator  CbCreator
	checkor  CbCheckor
}

func (c *CircuitBreakers) CreateBreaker(name string) {
	c.breakers.Store(name, c.creator(name))
}

func (c *CircuitBreakers) GetBreaker(name string) *CircuitBreaker {
	x, _ := c.breakers.LoadOrStore(name, c.creator(name))
	return x.(*CircuitBreaker)
}

func (c *CircuitBreakers) Clear() {
	if c.checkor != nil {
		c.breakers.Range(func(k, v interface{}) bool {
			if k != "" && !c.checkor(k.(string)) {
				log.Warn("delete circuit breaker: %s", k)
				c.breakers.Delete(k)
			}
			return true
		})
	}

}
