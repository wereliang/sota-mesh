package http

import (
	"fmt"

	"github.com/wereliang/sota-mesh/pkg/api"
)

type Handler interface {
	api.HTTPFilterManager
	api.DecoderFilterCallbacks
	Decode(api.StreamContext) error
	Encode(api.StreamContext) error
}

func NewHandler(r api.RouteConfigMatcher, c api.Connection) Handler {
	return &httpHandler{routeMatcher: r, connection: c}
}

type httpHandler struct {
	decodeFilters []api.StreamDecoderFilter
	encodeFilters []api.StreamEncoderFilter
	routeMatcher  api.RouteConfigMatcher
	connection    api.Connection
}

func (h *httpHandler) Route() api.RouteConfigMatcher {
	return h.routeMatcher
}

func (h *httpHandler) SetRoute(r api.RouteConfigMatcher) {
	h.routeMatcher = r
}

func (h *httpHandler) Connection() api.Connection {
	return h.connection
}

func (h *httpHandler) AddDecodeFilter(f api.StreamDecoderFilter) {
	f.SetDecoderFilterCallbacks(h)
	h.decodeFilters = append(h.decodeFilters, f)
}

func (h *httpHandler) AddEncodeFilter(f api.StreamEncoderFilter) {
	h.encodeFilters = append(h.encodeFilters, f)
}

func (h *httpHandler) Decode(ctx api.StreamContext) error {
	for _, f := range h.decodeFilters {
		if f.Decode(ctx) == api.Stop {
			return fmt.Errorf("decode error")
		}
	}
	return nil
}

func (h *httpHandler) Encode(ctx api.StreamContext) error {
	for _, f := range h.encodeFilters {
		if f.Encode(ctx) == api.Stop {
			return fmt.Errorf("encode error")
		}
	}
	return nil
}
