package config

type Filter interface {
	GetName() string
	GetTypedConfig() TypeConfig
}

type FilterImpl struct {
	Name        string         `json:"name"`
	TypedConfig TypeConfigImpl `json:"typed_config"`
}

func (l *FilterImpl) GetName() string {
	return l.Name
}

func (l *FilterImpl) GetTypedConfig() TypeConfig {
	if l.TypedConfig == nil {
		return nil
	}
	return l.TypedConfig
}
