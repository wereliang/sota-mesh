package config

type CidrRangeImpl struct {
	AddressPrefix string `json:"address_prefix"`
	PrefixLen     uint32 `json:"prefix_len"`
}

func (c *CidrRangeImpl) GetAddressPrefix() string {
	return c.AddressPrefix
}

func (c *CidrRangeImpl) GetPrefixLen() uint32 {
	return c.PrefixLen
}

type BindConfigImpl struct {
	SocketAddress *SocketAddress `json:"source_address"`
}

func (b *BindConfigImpl) GetSourceAddress() *SocketAddress {
	return b.SocketAddress
}
