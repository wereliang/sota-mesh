package main

import (
	_ "github.com/wereliang/sota-mesh/pkg/filter/envoy/http/router"
	_ "github.com/wereliang/sota-mesh/pkg/filter/envoy/listener/http_inspector"
	_ "github.com/wereliang/sota-mesh/pkg/filter/envoy/network/http_connection_manager"
	_ "github.com/wereliang/sota-mesh/pkg/filter/http/sota_router"
	_ "github.com/wereliang/sota-mesh/pkg/filter/network/simple_line"
	_ "github.com/wereliang/sota-mesh/pkg/filter/network/zmq"
)

func main() {
	Execute()
}
