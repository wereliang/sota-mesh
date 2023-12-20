package route

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"github.com/wereliang/sota-mesh/pkg/api"
	"github.com/wereliang/sota-mesh/pkg/config"
)

type requestHeader struct {
	*fasthttp.RequestHeader
	uri *fasthttp.URI
}

func (h *requestHeader) Get(key string) []byte {
	return h.RequestHeader.Peek(key)
}

func (h *requestHeader) Path() []byte {
	return h.uri.Path()
}

func (h *requestHeader) SetPath(path string) {
	h.uri.SetPath(path)
}

func newRequestHeader(request *fasthttp.Request) api.RequestHeader {
	return &requestHeader{RequestHeader: &request.Header, uri: request.URI()}
}

func TestHost(t *testing.T) {
	rc := &config.RouteConfigurationImpl{
		VirtualHost: []*config.VirtualHostImpl{
			{
				Name:    "test001",
				Domains: []string{"www.qq.com"},
				Routes: []*config.RouteImpl{
					{
						Match: &config.RouteMatchImpl{Path: "/foo"},
						Route: &config.RouteActionImpl{Cluster: "cluster_qq_foo"},
					},
					{
						Match: &config.RouteMatchImpl{Prefix: "/bar"},
						Route: &config.RouteActionImpl{Cluster: "cluster_qq_bar"},
					},
				},
			},
			{
				Name:    "test002",
				Domains: []string{"*.baidu.com"},
				Routes: []*config.RouteImpl{
					{
						Match: &config.RouteMatchImpl{Path: "/foo"},
						Route: &config.RouteActionImpl{Cluster: "cluster_baidu_foo"},
					},
				},
			},
			{
				Name:    "test003",
				Domains: []string{"www.ali.*"},
				Routes: []*config.RouteImpl{
					{
						Match: &config.RouteMatchImpl{Path: "/foo"},
						Route: &config.RouteActionImpl{Cluster: "cluster_ali_foo"},
					},
				},
			},
			{
				Name:    "test004",
				Domains: []string{"*"},
				Routes: []*config.RouteImpl{
					{
						Match: &config.RouteMatchImpl{Prefix: "/"},
						Route: &config.RouteActionImpl{Cluster: "cluster_wildcar"},
					},
				},
			},
		},
	}

	tests := []struct {
		host    string
		path    string
		isnil   bool
		cluster string
	}{
		{
			"www.qq.com",
			"/foo",
			false,
			"cluster_qq_foo",
		},
		{
			"www.qq.com",
			"/bar/xxx",
			false,
			"cluster_qq_bar",
		},
		{
			"xxx.baidu.com",
			"/foo",
			false,
			"cluster_baidu_foo",
		},
		{
			"yyy.baidu.com",
			"/abc",
			true,
			"",
		},
		{
			"www.ali.cn",
			"/foo",
			false,
			"cluster_ali_foo",
		},
		{
			"www.xxx.com",
			"/foo",
			false,
			"cluster_wildcar",
		},
	}

	for _, item := range tests {
		request := &fasthttp.Request{}
		header := &request.Header
		header.SetRequestURI(item.path)
		header.SetHost(item.host)

		matcher := NewRouterMatcher(rc)
		entry := matcher.Match(newRequestHeader(request))
		if item.isnil {
			assert.Nil(t, entry)
		} else {
			assert.Equal(t, entry.GetClusterName(), item.cluster)
		}
	}
}
