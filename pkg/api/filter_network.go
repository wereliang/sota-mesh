package api

import (
	"bufio"
	"bytes"
)

// ReadFilter for network filter
type ReadFilter interface {
	// OnData is called everytime bytes is read from the connection
	OnData(*bytes.Buffer) FilterStatus
	// OnNewConnection is called on new connection is created
	OnNewConnection() FilterStatus
}

// WriteFilter for network filter
type WriteFilter interface {
	// OnWrite is called before data write to raw connection
	OnWrite(*bufio.Writer) FilterStatus
}

// FilterManager manager network filter
type FilterManager interface {
	AddReadFilter(ReadFilter)
	AddWriteFilter(WriteFilter)
}

type NetworkFilterCreator func(FilterManager, ConnectionCallbacks) error

// NetworkFactory for network factory
type NetworkFactory interface {
	Factory
	CreateFilterFactory(interface{}, FactoryContext) NetworkFilterCreator
}
