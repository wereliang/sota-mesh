package filter

import (
	"fmt"

	"github.com/wereliang/sota-mesh/pkg/api"
	"github.com/wereliang/sota-mesh/pkg/config"
	"github.com/wereliang/sota-mesh/pkg/log"
)

var (
	ListenerFilterFactory = newRegistFactory()
	NetworkFilterFactory  = newRegistFactory()
	HTTPFilterFactory     = newRegistFactory()
)

func GetListenerFactory(a config.TypeConfig, name string) (api.ListernerFactory, interface{}) {
	factory, cfg := getFactory(a, name, ListenerFilterFactory)
	if factory != nil {
		return factory.(api.ListernerFactory), cfg
	}
	return nil, cfg
}

func GetNetworkFactory(a config.TypeConfig, name string) (api.NetworkFactory, interface{}) {
	factory, cfg := getFactory(a, name, NetworkFilterFactory)
	if factory != nil {
		return factory.(api.NetworkFactory), cfg
	}
	return nil, cfg
}

func GetHTTPFactory(a config.TypeConfig, name string) (api.HTTPFactory, interface{}) {
	factory, pb := getFactory(a, name, HTTPFilterFactory)
	if factory != nil {
		return factory.(api.HTTPFactory), pb
	}
	return nil, pb
}

func getFactory(a config.TypeConfig, name string, r *registFactory) (api.Factory, interface{}) {
	var (
		factory api.Factory
		cfg     interface{}
	)
	if a == nil {
		factory = r.GetFactoryByName(name)
	} else {
		if factory = r.GetFactoryByType(a.Type()); factory == nil {
			return nil, nil
		}

		cfg = factory.CreateEmptyConfigProto()
		err := config.UnmarshalConfigProto(a, cfg)
		if err != nil {
			panic(err)
		}
		log.Debug("cfg:%#v", cfg)
	}
	return factory, cfg
}

type registFactory struct {
	namedFactorys map[string]api.Factory
	typedFactorys map[string]api.Factory
}

func newRegistFactory() *registFactory {
	return &registFactory{
		namedFactorys: make(map[string]api.Factory),
		typedFactorys: make(map[string]api.Factory),
	}
}

func (f *registFactory) Regist(factory api.Factory) {
	name := factory.Name()
	if name == "" {
		panic("name can't be empty")
	}
	if _, ok := f.namedFactorys[name]; ok {
		panic("factory name exist")
	}
	log.Debug("regist factory: %s", name)
	f.namedFactorys[name] = factory

	if cfg := factory.CreateEmptyConfigProto(); cfg != nil {
		typed, err := config.GetConfigProtoType(cfg)
		if err != nil || typed == "" {
			panic(fmt.Errorf("get config proto type error. %s [%s]", err, typed))
		}
		if _, ok := f.typedFactorys[typed]; ok {
			panic("factory typed exist")
		}
		f.typedFactorys[typed] = factory
	}
}

func (f *registFactory) GetFactoryByName(name string) api.Factory {
	if f, ok := f.namedFactorys[name]; ok {
		return f
	}
	return nil
}

func (f *registFactory) GetFactoryByType(typed string) api.Factory {
	if f, ok := f.typedFactorys[typed]; ok {
		return f
	}
	return nil
}

func (f *registFactory) RegisteredNames() []string {
	var names []string
	for n := range f.namedFactorys {
		names = append(names, n)
	}
	return names
}

func (f *registFactory) RegisteredTypes() []string {
	var types []string
	for n := range f.typedFactorys {
		types = append(types, n)
	}
	return types
}
