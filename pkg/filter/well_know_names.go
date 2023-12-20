package filter

const (
	Listener_TlsInspector         = "envoy.filters.listener.tls_inspector"
	Listener_OriginalDst          = "envoy.filters.listener.original_dst"
	Listener_HttpInspector        = "envoy.filters.listener.http_inspector"
	Network_Echo                  = "envoy.filters.network.echo"
	Network_HttpConnectionManager = "envoy.filters.network.http_connection_manager"
	HTTP_Router                   = "envoy.filters.http.router"
)

var well_know_names = map[string]struct{}{
	Listener_TlsInspector:         {},
	Listener_OriginalDst:          {},
	Listener_HttpInspector:        {},
	Network_HttpConnectionManager: {},
	HTTP_Router:                   {},
}

func IsWellknowName(name string) bool {
	_, ok := well_know_names[name]
	return ok
}
