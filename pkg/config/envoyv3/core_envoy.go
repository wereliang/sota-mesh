package envoyv3

import envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"

type CidrRangeEnvoy struct {
	*envoy_config_core_v3.CidrRange
}

func (c *CidrRangeEnvoy) GetAddressPrefix() string {
	return c.CidrRange.GetAddressPrefix()
}

func (c *CidrRangeEnvoy) GetPrefixLen() uint32 {
	return c.CidrRange.GetPrefixLen().GetValue()
}
