package envoyv3

import (
	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/access_loggers/file/v3"
	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/cors/v3"
	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/fault/v3"
	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/grpc_stats/v3"
	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/router/v3"
	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/wasm/v3"
	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/listener/http_inspector/v3"
	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/tcp_proxy/v3"

	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/request_id/uuid/v3"
	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/upstreams/http/v3"

	_ "github.com/cncf/xds/go/udpa/type/v1"

	_ "istio.io/api/envoy/config/filter/http/alpn/v2alpha1"
	_ "istio.io/api/envoy/config/filter/network/metadata_exchange"

	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
)
