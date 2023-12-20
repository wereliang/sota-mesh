package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/wereliang/sota-mesh/pkg/log"
)

type CfgType string

const (
	CTypeSota  CfgType = "sota"
	CTypeEnvoy CfgType = "envoy"
)

type SotaConfig interface {
	GetBootstrap() Bootstrap
	Dump()
}

type SotaConfigImpl struct {
	bootstrap Bootstrap
}

func (s *SotaConfigImpl) GetBootstrap() Bootstrap {
	return s.bootstrap
}

func (s *SotaConfigImpl) Dump() {
	bs := s.bootstrap.(*BootstrapImpl)
	data, err := json.Marshal(bs)
	if err != nil {
		panic(err)
	}
	log.Debug(string(data))
	return
}

func NewConfig(path string) (SotaConfig, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	// TODO: check is yaml
	content, err = yaml.YAMLToJSON(content)
	if err != nil {
		return nil, err
	}

	bs := &BootstrapImpl{}
	err = json.Unmarshal(content, bs)
	if err != nil {
		return nil, err
	}
	return &SotaConfigImpl{bs}, nil
}
