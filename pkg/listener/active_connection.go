package listener

import (
	"bytes"
	"fmt"
	"io"

	"github.com/wereliang/sota-mesh/pkg/api"
	"github.com/wereliang/sota-mesh/pkg/log"
)

func NewActiveConnection(conn api.Connection) api.ActiveConnection {
	return &activeConnection{
		Connection: conn,
	}
}

type activeConnection struct {
	api.Connection
	rfs []api.ReadFilter
	wfs []api.WriteFilter
}

func (ac *activeConnection) AddReadFilter(f api.ReadFilter) {
	// TODO: 现阶段只支持一个network filter
	if len(ac.rfs) > 0 {
		panic("network filter just support only one")
	}
	ac.rfs = append(ac.rfs, f)
}

func (ac *activeConnection) AddWriteFilter(f api.WriteFilter) {
	ac.wfs = append(ac.wfs, f)
}

func (ac *activeConnection) Close() error {
	return ac.Connection.Close()
}

func (ac *activeConnection) close() {
	log.Trace("connection close")
	ac.Connection.Close()
}

func (ac *activeConnection) OnLoop() {
	defer ac.close()

	if !ac.onNewConnection() {
		return
	}

	bytesBuf := bytes.NewBuffer(make([]byte, 0, 1024))
	for {
		bs := make([]byte, 1024)
		n, err := ac.Read(bs)
		if err != nil {
			if err != io.EOF {
				bytesBuf.Reset()
				log.Error("read error: %s", err)
			}
			// 下一层收到空的buffer需要清理资源
			ac.onData(bytesBuf)
			break
		}

		_, e := bytesBuf.Write(bs[:n])
		if e != nil {
			fmt.Println("bytes buf write error", e)
			return
		}
		if !ac.onData(bytesBuf) {
			break
		}
	}
}

func (ac *activeConnection) onNewConnection() bool {
	for _, f := range ac.rfs {
		if f.OnNewConnection() != api.Continue {
			return false
		}
	}
	return true
}

func (ac *activeConnection) onData(buf *bytes.Buffer) bool {
	for _, f := range ac.rfs {
		if f.OnData(buf) != api.Continue {
			return false
		}
	}
	return true
}
