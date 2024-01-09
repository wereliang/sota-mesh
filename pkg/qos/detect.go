package qos

import (
	"fmt"
	"net"
	"time"

	"github.com/wereliang/sota-mesh/pkg/log"
)

// TCPDetect : TCP连通性探测
func TCPDetect(addr string) bool {
	conn, err := net.DialTimeout("tcp", addr, 100*time.Millisecond)
	if err != nil {
		log.Error("Tcp Connect(%s) Error:%s", addr, err.Error())
		return false
	}
	tcpConn, _ := conn.(*net.TCPConn)
	tcpConn.SetLinger(0)
	tcpConn.Close()
	fmt.Printf("Tcp Connect(%s) ok\n", addr)
	return true
}
