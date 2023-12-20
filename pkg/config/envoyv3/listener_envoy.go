package envoyv3

import (
	"net"

	envoy_config_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	"github.com/wereliang/sota-mesh/pkg/config"
)

type ListenerEnvoy struct {
	*envoy_config_listener_v3.Listener
}

func (l *ListenerEnvoy) GetAddress() net.Addr {
	addr, err := ToNetAddr(l.Address)
	if err != nil {
		panic(err)
	}
	return addr
}

func (l *ListenerEnvoy) GetListenerFilters() []config.Filter {
	var filters []config.Filter
	for _, el := range l.ListenerFilters {
		filters = append(filters, &ListernerFilterEnvoy{el})
	}
	return filters
}

func (l *ListenerEnvoy) GetFilterChains() []config.FilterChain {
	var chains []config.FilterChain
	for _, fc := range l.FilterChains {
		chains = append(chains, &FilterChainEnvoy{fc})
	}
	return chains
}

func (l *ListenerEnvoy) GetDefaultFilterChain() config.FilterChain {
	if x := l.Listener.GetDefaultFilterChain(); x != nil {
		return &FilterChainEnvoy{x}
	}
	return nil
}

func (l *ListenerEnvoy) GetUseOriginalDst() bool {
	return l.Listener.GetUseOriginalDst().GetValue()
}

func (l *ListenerEnvoy) GetBindToPort() bool {
	// default is true
	if x := l.Listener.GetBindToPort(); x != nil {
		return x.GetValue()
	}
	return true
}

// ///////////////////////
type ListernerFilterEnvoy struct {
	*envoy_config_listener_v3.ListenerFilter
}

func (l *ListernerFilterEnvoy) GetTypedConfig() config.TypeConfig {
	if any := l.ListenerFilter.GetTypedConfig(); any != nil {
		return &TypeConfigEnvoy{any}
	}
	return nil
}

// ////////////////////
type FilterChainEnvoy struct {
	*envoy_config_listener_v3.FilterChain
}

func (f *FilterChainEnvoy) GetFilterChainMatch() config.FilterChainMatch {
	if x := f.FilterChain.GetFilterChainMatch(); x != nil {
		return &FilterChanMatchEnvoy{x}
	}
	return nil
}

func (f *FilterChainEnvoy) GetFilters() []config.Filter {
	var filters []config.Filter
	for _, filter := range f.FilterChain.GetFilters() {
		filters = append(filters, &FilterEnvoy{filter})
	}
	return filters
}

func (f *FilterChainEnvoy) SetName(s string) {
	f.FilterChain.Name = s
}

// //////////////////////////////
type FilterChanMatchEnvoy struct {
	*envoy_config_listener_v3.FilterChainMatch
}

func (f *FilterChanMatchEnvoy) GetPrefixRanges() []config.CidrRange {
	var cranges []config.CidrRange
	for _, r := range f.FilterChainMatch.GetPrefixRanges() {
		cranges = append(cranges, &CidrRangeEnvoy{r})
	}
	return cranges
}

func (f *FilterChanMatchEnvoy) GetDirectSourcePrefixRanges() []config.CidrRange {
	var cranges []config.CidrRange
	for _, r := range f.FilterChainMatch.GetDirectSourcePrefixRanges() {
		cranges = append(cranges, &CidrRangeEnvoy{r})
	}
	return cranges
}

func (f *FilterChanMatchEnvoy) GetSourceType() config.FilterChainMatch_ConnectionSourceType {
	return config.FilterChainMatch_ConnectionSourceType(f.FilterChainMatch.GetSourceType())
}

func (f *FilterChanMatchEnvoy) GetSourcePrefixRanges() []config.CidrRange {
	var cranges []config.CidrRange
	for _, r := range f.FilterChainMatch.GetSourcePrefixRanges() {
		cranges = append(cranges, &CidrRangeEnvoy{r})
	}
	return cranges
}

func (f *FilterChanMatchEnvoy) GetDestinationPort() uint32 {
	return f.FilterChainMatch.GetDestinationPort().GetValue()
}

// //////////////////////
type FilterEnvoy struct {
	*envoy_config_listener_v3.Filter
}

func (l *FilterEnvoy) GetTypedConfig() config.TypeConfig {
	if any := l.Filter.GetTypedConfig(); any != nil {
		return &TypeConfigEnvoy{any}
	}
	return nil
}
