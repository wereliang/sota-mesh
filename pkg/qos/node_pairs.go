package qos

import (
	"github.com/wereliang/sota-mesh/pkg/config"
)

type NodePairs struct {
	edps  config.LbEndpointSet
	index map[string]int
}

func newNodeParis(edps config.LbEndpointSet) NodePairs {
	np := NodePairs{}
	np.Refresh(edps)
	return np
}

func (np *NodePairs) String() string {
	var res string
	for _, e := range np.edps {
		res += hashEndpoint(e) + " "
	}
	return res
}

func (np *NodePairs) Refresh(edps config.LbEndpointSet) {
	np.edps = edps
	np.index = make(map[string]int)
	for i, e := range np.edps {
		np.index[hashEndpoint(e)] = i
	}
}

func (np *NodePairs) Exist(e config.LbEndpoint) bool {
	return np.ExistByStr(hashEndpoint(e))
}

func (np *NodePairs) ExistByStr(str string) bool {
	_, ok := np.index[str]
	return ok
}

func (np *NodePairs) Add(e config.LbEndpoint) {
	if !np.Exist(e) {
		np.edps = append(np.edps, e)
		np.index[hashEndpoint(e)] = len(np.edps) - 1
	}
}

func (np *NodePairs) Del(e config.LbEndpoint) (config.LbEndpoint, bool) {
	return np.DelByStr(hashEndpoint(e))
}

func (np *NodePairs) DelByStr(str string) (config.LbEndpoint, bool) {
	if i, ok := np.index[str]; ok {
		e := np.edps[i]
		edps := append(np.edps[:i], np.edps[i+1:]...)
		np.Refresh(edps)
		return e, true
	}
	return nil, false
}

func hashEndpoint(e config.LbEndpoint) string {
	return e.GetEndpoint().GetAddress().String()
}
