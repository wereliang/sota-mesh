package config

import (
	"net"
)

type ListenerImpl struct {
	Name               string             `json:"name"`
	Address            *Address           `json:"address"`
	ListenerFilters    []*FilterImpl      `json:"listener_filters"`
	FilterChains       []*FilterChainImpl `json:"filter_chains"`
	DefaultFilterChain *FilterChainImpl   `json:"default_filter_chain"`
	UseOriginalDst     bool               `json:"use_original_dst"`
	BindToPort         bool               `json:"bind_to_port" default:"true"`
}

func (l *ListenerImpl) GetName() string {
	return l.Name
}

func (l *ListenerImpl) GetAddress() net.Addr {
	return l.Address
}

func (l *ListenerImpl) GetListenerFilters() []Filter {
	var filters []Filter
	for _, lf := range l.ListenerFilters {
		filters = append(filters, lf)
	}
	return filters
}

func (l *ListenerImpl) GetFilterChains() []FilterChain {
	var chains []FilterChain
	for _, fc := range l.FilterChains {
		chains = append(chains, fc)
	}
	return chains
}

func (l *ListenerImpl) GetDefaultFilterChain() FilterChain {
	if l.DefaultFilterChain == nil {
		return nil
	}
	return l.DefaultFilterChain
}

func (l *ListenerImpl) GetUseOriginalDst() bool {
	return l.UseOriginalDst
}

func (l *ListenerImpl) GetBindToPort() bool {
	return l.BindToPort
}

type FilterChainImpl struct {
	Name             string                `json:"name"`
	FilterChainMatch *FilterChainMatchImpl `json:"filter_chain_match"`
	Filters          []*FilterImpl         `json:"filters"`
}

func (f *FilterChainImpl) GetName() string {
	return f.Name
}

func (f *FilterChainImpl) SetName(s string) {
	f.Name = s
}

func (f *FilterChainImpl) GetFilterChainMatch() FilterChainMatch {
	if f.FilterChainMatch == nil {
		return nil
	}
	return f.FilterChainMatch
}

func (f *FilterChainImpl) GetFilters() []Filter {
	var filters []Filter
	for _, filter := range f.Filters {
		filters = append(filters, filter)
	}
	return filters
}

type FilterChainMatchImpl struct {
	DestinationPort          uint32                                `json:"destination_port"`
	PrefixRanges             []*CidrRangeImpl                      `json:"prefix_ranges"`
	ServerNames              []string                              `json:"server_names"`
	ApplicationProtocols     []string                              `json:"application_protocols"`
	DirectSourcePrefixRanges []*CidrRangeImpl                      `json:"direct_source_prefix_ranges"`
	SourceType               FilterChainMatch_ConnectionSourceType `json:"source_type"`
	SourcePrefixRanges       []*CidrRangeImpl                      `json:"source_prefix_ranges"`
	SourcePorts              []uint32                              `json:"source_ports"`
	TransportProtocol        string                                `json:"transport_protocol"`
}

func (f *FilterChainMatchImpl) GetPrefixRanges() []CidrRange {
	var cidranges []CidrRange
	for _, cr := range f.PrefixRanges {
		cidranges = append(cidranges, cr)
	}
	return cidranges
}

func (f *FilterChainMatchImpl) GetServerNames() []string {
	return f.ServerNames
}

func (f *FilterChainMatchImpl) GetApplicationProtocols() []string {
	return f.ApplicationProtocols
}

func (f *FilterChainMatchImpl) GetDirectSourcePrefixRanges() []CidrRange {
	var cidranges []CidrRange
	for _, cr := range f.DirectSourcePrefixRanges {
		cidranges = append(cidranges, cr)
	}
	return cidranges
}

func (f *FilterChainMatchImpl) GetSourceType() FilterChainMatch_ConnectionSourceType {
	return f.SourceType
}

func (f *FilterChainMatchImpl) GetSourcePrefixRanges() []CidrRange {
	var cidranges []CidrRange
	for _, cr := range f.SourcePrefixRanges {
		cidranges = append(cidranges, cr)
	}
	return cidranges
}

func (f *FilterChainMatchImpl) GetSourcePorts() []uint32 {
	return f.SourcePorts
}

func (f *FilterChainMatchImpl) GetDestinationPort() uint32 {
	return f.DestinationPort
}

func (f *FilterChainMatchImpl) GetTransportProtocol() string {
	return f.TransportProtocol
}
