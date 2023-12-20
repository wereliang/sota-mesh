package http

import (
	"context"

	"github.com/valyala/fasthttp"
	"github.com/wereliang/sota-mesh/pkg/api"
)

func NewStreamContext(context context.Context,
	req *fasthttp.Request, rsp *fasthttp.Response) api.StreamContext {
	return &streamContext{
		context:  context,
		request:  newRequest(req),
		response: newResponse(rsp),
	}
}

type streamContext struct {
	context  context.Context
	request  api.Request
	response api.Response
}

func (sc *streamContext) Context() context.Context {
	return sc.context
}

func (sc *streamContext) Request() api.Request {
	return sc.request
}

func (sc *streamContext) Response() api.Response {
	return sc.response
}

func newRequest(r *fasthttp.Request) api.Request {
	return &request{
		Request: r,
		header:  newRequestHeader(r),
		body:    newRequestBody(r),
	}
}

type request struct {
	*fasthttp.Request
	header api.RequestHeader
	body   api.Body
}

func (r *request) Header() api.RequestHeader {
	return r.header
}

func (r *request) Body() api.Body {
	return r.body
}

func (r *request) Raw() interface{} {
	return r.Request
}

func newResponse(r *fasthttp.Response) api.Response {
	return &response{
		Response: r,
		header:   newResponseHeader(r),
		body:     newResponseBody(r),
	}
}

type response struct {
	*fasthttp.Response
	header api.ResponseHeader
	body   api.Body
}

func (r *response) Header() api.ResponseHeader {
	return r.header
}

func (r *response) Body() api.Body {
	return r.body
}

func (r *response) Raw() interface{} {
	return r.Response
}

func newRequestHeader(request *fasthttp.Request) api.RequestHeader {
	return &requestHeader{RequestHeader: &request.Header, uri: request.URI()}
}

func newResponseHeader(response *fasthttp.Response) api.ResponseHeader {
	return &responseHeader{&response.Header}
}

func newRequestBody(request *fasthttp.Request) api.Body {
	return &requestBody{Request: request}
}

func newResponseBody(response *fasthttp.Response) api.Body {
	return &responseBody{Response: response}
}

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

type responseHeader struct {
	*fasthttp.ResponseHeader
}

func (h *responseHeader) Get(key string) []byte {
	return h.ResponseHeader.Peek(key)
}

type requestBody struct {
	*fasthttp.Request
}

type responseBody struct {
	*fasthttp.Response
}
