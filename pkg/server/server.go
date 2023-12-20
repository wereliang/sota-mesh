package server

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/wereliang/sota-mesh/pkg/api"
	"github.com/wereliang/sota-mesh/pkg/config/envoyv3"
	"github.com/wereliang/sota-mesh/pkg/log"

	"github.com/wereliang/sota-mesh/pkg/cluster"
	"github.com/wereliang/sota-mesh/pkg/config"
	"github.com/wereliang/sota-mesh/pkg/listener"
)

type Sota struct {
	lm  api.ListenerManager
	cm  api.ClusterManager
	rm  api.RouteConfigManager
	cfg config.SotaConfig
}

func NewSota(path string, ctype config.CfgType) (*Sota, error) {
	var err error
	var cfg config.SotaConfig
	switch ctype {
	case config.CTypeSota:
		cfg, err = config.NewConfig(path)
	case config.CTypeEnvoy:
		cfg, err = envoyv3.NewEnvoyConfig(path)
	default:
		panic(fmt.Sprintf("invalid ctype: %s", ctype))
	}
	if err != nil {
		return nil, err
	}
	return &Sota{cfg: cfg}, nil
}

func (s *Sota) Start() error {
	s.init()
	s.waitSignal()
	s.stop()
	return nil
}

func (s *Sota) init() error {
	s.initLogger()
	s.cfg.Dump()

	var err error
	if err = s.initAdmin(); err != nil {
		return err
	}
	if err = s.initCluster(); err != nil {
		return err
	}
	if err = s.initRouteConfig(); err != nil {
		return err
	}
	if err = s.initListener(); err != nil {
		return err
	}
	return nil
}

func (s *Sota) initLogger() {
}

func (s *Sota) initAdmin() error {
	return nil
}

func (s *Sota) initCluster() error {
	var err error
	if s.cm, err = cluster.NewClusterManager(
		s.cfg.GetBootstrap().GetStaticResources().GetClusters()); err != nil {
		return err
	}
	return nil
}

func (s *Sota) initRouteConfig() error {
	return nil
}

func (s *Sota) initListener() error {
	var err error
	if s.lm, err = listener.NewListenerManager(s); err != nil {
		return err
	}

	for _, l := range s.cfg.GetBootstrap().GetStaticResources().GetListeners() {
		err = s.lm.AddOrUpdateListener(config.LISTENER_STATIC, l)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Sota) waitSignal() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	x := <-c
	log.Debug("catch signal %v", x)
}

func (s *Sota) stop() error {
	return nil
}

func (s *Sota) ClusterManager() api.ClusterManager {
	return s.cm
}

func (s *Sota) RouteConfigManager() api.RouteConfigManager {
	return s.rm
}

func (s *Sota) ListenerManager() api.ListenerManager {
	return s.lm
}
