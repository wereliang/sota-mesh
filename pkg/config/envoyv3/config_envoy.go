package envoyv3

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	envoy_config_bootstrap_v3 "github.com/envoyproxy/go-control-plane/envoy/config/bootstrap/v3"
	"github.com/ghodss/yaml"

	"github.com/golang/protobuf/jsonpb"
	"github.com/wereliang/sota-mesh/pkg/config"
	"github.com/wereliang/sota-mesh/pkg/log"
)

type SotaConfigEnvoy struct {
	*envoy_config_bootstrap_v3.Bootstrap
}

func (s *SotaConfigEnvoy) GetBootstrap() config.Bootstrap {
	return &BootstrapEnvoy{s.Bootstrap}
}

func (s *SotaConfigEnvoy) Dump() {
	fmt.Println(s.Bootstrap.String())
}

func NewEnvoyConfig(path string) (config.SotaConfig, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	bootrap := &envoy_config_bootstrap_v3.Bootstrap{}

	if yamlFormat(path) {
		content, err = yaml.YAMLToJSON(content)
		if err != nil {
			return nil, err
		}
	}

	if err = jsonpb.UnmarshalString(string(content), bootrap); err != nil {
		log.Error("jsonpb.UnmarshalString error %s", err)
		return nil, err
	}

	if err = bootrap.Validate(); err != nil {
		return nil, err
	}
	return &SotaConfigEnvoy{bootrap}, nil
}

func yamlFormat(path string) bool {
	ext := filepath.Ext(path)
	if ext == ".yaml" || ext == ".yml" {
		return true
	}
	return false
}
