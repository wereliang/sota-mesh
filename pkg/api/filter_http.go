package api

// StreamDecoderFilter for http stream filter
type StreamDecoderFilter interface {
	// Decode call before request
	Decode(StreamContext) FilterStatus
	SetDecoderFilterCallbacks(DecoderFilterCallbacks)
}

// DecoderFilterCallbacks
type DecoderFilterCallbacks interface {
	Connection() Connection
	Route() RouteConfigMatcher
	SetRoute(RouteConfigMatcher)
}

// StreamEncoderFilter for http stream filter
type StreamEncoderFilter interface {
	// Encode call after request
	Encode(StreamContext) FilterStatus
}

// HTTPFilterManager manager http filter
type HTTPFilterManager interface {
	AddDecodeFilter(StreamDecoderFilter)
	AddEncodeFilter(StreamEncoderFilter)
}

type HTTPFilterCreator func(HTTPFilterManager)

// HTTPFactory for http factory
type HTTPFactory interface {
	Factory
	CreateFilterFactory(interface{}, FactoryContext) HTTPFilterCreator
}
