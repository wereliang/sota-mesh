package config

import (
	"fmt"
	"strings"
)

type SocketAddress struct {
	Protocol  string `json:"protocol"`
	Address   string `json:"address"`
	PortValue uint32 `json:"port_value"`
}

type Address struct {
	SocketAddress SocketAddress `json:"socket_address"`
}

// Network name of the network (for example, "tcp", "udp")
func (addr *Address) Network() string {
	return strings.ToLower(addr.SocketAddress.Protocol)
}

// String string form of address (for example, "192.0.2.1:25", "[2001:db8::1]:80")
func (addr *Address) String() string {
	return fmt.Sprintf("%s:%d", addr.SocketAddress.Address, addr.SocketAddress.PortValue)
}

// CidrRange cidr range
type CidrRange interface {
	GetAddressPrefix() string
	GetPrefixLen() uint32
}

// BindConfig
type BindConfig interface {
	GetSourceAddress() *SocketAddress
}
