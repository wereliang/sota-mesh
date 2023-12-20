package dispatch

import (
	"bufio"
	"bytes"
	"fmt"
	"io"

	"github.com/wereliang/sota-mesh/pkg/api"
)

type Dispatcher interface {
	Dispatch(*bytes.Buffer) error
}

type DispatcherImpl struct {
	bufChan   chan *bytes.Buffer
	endChan   chan struct{}
	closeChan chan struct{}
	Reader    *bufio.Reader
	Writer    *bufio.Writer
	Conn      api.ConnectionCallbacks
}

func (dp *DispatcherImpl) Dispatch(buffer *bytes.Buffer) error {
	select {
	case <-dp.closeChan:
		return fmt.Errorf("stream server close")
	default:
		dp.bufChan <- buffer
		<-dp.endChan
	}
	return nil
}

func (dp *DispatcherImpl) Read(dst []byte) (n int, err error) {
	buf, ok := <-dp.bufChan
	if !ok {
		err = io.EOF
	} else {
		n, err = buf.Read(dst)
	}
	dp.endChan <- struct{}{}
	return n, err
}

func (dp *DispatcherImpl) Write(p []byte) (n int, err error) {
	return dp.Conn.Write(p)
}

func (dp *DispatcherImpl) Close() {
	close(dp.closeChan)
	dp.Conn.Close()
}

func NewDispatcher(conn api.Connection) (*DispatcherImpl, error) {
	dp := &DispatcherImpl{
		bufChan:   make(chan *bytes.Buffer),
		endChan:   make(chan struct{}),
		closeChan: make(chan struct{}),
		Conn:      conn,
	}
	dp.Reader = bufio.NewReader(dp)
	dp.Writer = bufio.NewWriter(dp)
	return dp, nil
}
