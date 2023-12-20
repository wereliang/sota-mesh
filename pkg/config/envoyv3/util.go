package envoyv3

import (
	"fmt"
	"net"
	"strings"

	envoy_config_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
)

func ToNetAddr(addr *envoy_config_v3.Address) (net.Addr, error) {
	return toNetAddr(addr, func(saddr *envoy_config_v3.SocketAddress) net.Addr {
		return &net.TCPAddr{IP: net.ParseIP(saddr.GetAddress()), Port: int(saddr.GetPortValue())}
	})
}

func toNetAddr(
	addr *envoy_config_v3.Address,
	fn func(*envoy_config_v3.SocketAddress) net.Addr) (net.Addr, error) {

	var naddr net.Addr
	saddr := addr.GetSocketAddress()
	switch strings.ToLower(saddr.GetProtocol().String()) {
	case "tcp":
		naddr = fn(saddr)
	default:
		return nil, fmt.Errorf("not support protocol")
	}
	return naddr, nil
}

type TypeConfigEnvoy struct {
	*any.Any
}

func (tc *TypeConfigEnvoy) Type() string {
	return tc.GetTypeUrl()
}

func (tc *TypeConfigEnvoy) Unmarshal(obj interface{}) error {
	if v, ok := obj.(proto.Message); !ok {
		panic("obj is not proto message")
	} else {
		return ptypes.UnmarshalAny(tc.Any, v)
	}
}
