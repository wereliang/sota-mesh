package network

import (
	"bufio"
	"net"

	"github.com/wereliang/sota-mesh/pkg/api"
)

func newConn(c net.Conn) api.Connection {
	return &connection{Conn: c, ctx: newConnectionContext(c), reader: bufio.NewReader(c)}
}

func newConnSize(c net.Conn, size int) api.Connection {
	return &connection{Conn: c, ctx: newConnectionContext(c), reader: bufio.NewReaderSize(c, size)}
}

type connection struct {
	net.Conn
	ctx    api.ConnectionContext
	reader *bufio.Reader
}

func (c *connection) Context() api.ConnectionContext {
	return c.ctx
}

func (c *connection) Raw() net.Conn {
	return c.Conn
}

func (c *connection) Peek(n int) ([]byte, error) {
	return c.reader.Peek(n)
}

func (c *connection) Read(p []byte) (int, error) {
	return c.reader.Read(p)
}

func newConnectionContext(c net.Conn) api.ConnectionContext {
	cs := &ConnectionContextImpl{}
	remote := c.RemoteAddr().(*net.TCPAddr)
	cs.SourceIP = remote.IP
	cs.SourcePort = uint32(remote.Port)
	local := c.LocalAddr().(*net.TCPAddr)
	cs.DestinationIP = local.IP
	cs.DestinationPort = uint32(local.Port)
	return cs
}

type ConnectionContextImpl struct {
	DestinationPort      uint32
	DestinationIP        net.IP
	ServerName           string
	TransportProtocol    string
	ApplicationProtocol  string
	DirectSourceIP       net.IP
	SourceType           int32
	SourceIP             net.IP
	SourcePort           uint32
	localAddressRestored bool
}

func (cs *ConnectionContextImpl) GetDestinationPort() uint32 {
	return cs.DestinationPort
}

func (cs *ConnectionContextImpl) GetDestinationIP() net.IP {
	return cs.DestinationIP
}

func (cs *ConnectionContextImpl) GetServerName() string {
	return cs.ServerName
}

func (cs *ConnectionContextImpl) GetTransportProtocol() string {
	return cs.TransportProtocol
}

func (cs *ConnectionContextImpl) GetApplicationProtocol() string {
	return cs.ApplicationProtocol
}

func (cs *ConnectionContextImpl) GetDirectSourceIP() net.IP {
	return cs.DirectSourceIP
}

func (cs *ConnectionContextImpl) GetSourceType() int32 {
	return cs.SourceType
}

func (cs *ConnectionContextImpl) GetSourceIP() net.IP {
	return cs.SourceIP
}

func (cs *ConnectionContextImpl) GetSourcePort() uint32 {
	return cs.SourcePort
}

func (cs *ConnectionContextImpl) SetOriginalDestination(ip net.IP, port uint32) {
	cs.DestinationIP = ip
	cs.DestinationPort = port
	cs.localAddressRestored = true
}

func (cs *ConnectionContextImpl) SetApplicationProtocol(s string) {
	cs.ApplicationProtocol = s
}

func (cs *ConnectionContextImpl) LocalAddressRestored() bool {
	return cs.localAddressRestored
}
