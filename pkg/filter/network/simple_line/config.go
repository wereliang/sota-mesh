package simpleline

type SimpleConfig struct {
	Protocol string `json:"protocol"`
}

func (c *SimpleConfig) Type() string {
	return "type.sota.com/sota.filters.network.simple"
}
