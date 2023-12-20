package api

import "context"

// HeaderMap is some header action
type HeaderMap interface {
	Del(key string)
	Add(key, value string)
	Set(key, value string)
	Get(key string) []byte
}

// StreamContext http stream context
type StreamContext interface {
	Context() context.Context
	Request() Request
	Response() Response
}

// Request http request
type Request interface {
	Header() RequestHeader
	Body() Body
	SetHost(string)
	Raw() interface{}
}

// Response http response
type Response interface {
	Header() ResponseHeader
	Body() Body
	Raw() interface{}
}

// RequestHeader http request header
type RequestHeader interface {
	HeaderMap
	Method() []byte
	SetMethod(method string)
	Host() []byte
	SetHost(host string)
	RequestURI() []byte
	SetRequestURI(requestURI string)
	Path() []byte
	SetPath(path string)
}

// ResponseHeader http response header
type ResponseHeader interface {
	HeaderMap
	StatusCode() int
	SetStatusCode(statusCode int)
}

// Body http body action
type Body interface {
	AppendBody([]byte)
	SetBody([]byte)
}
