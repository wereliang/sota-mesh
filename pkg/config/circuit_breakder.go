package config

// 通过限制某个客户端对目标服务的连接数、访问 请求数等，避免对一个服务的过量访问，
// 如果超过配置的阈值，则快 速断路请求。
// 还会限制重试次数，避免重试次数过多导致系统压力变大并加剧故障的传播
type CircuitBreakers interface {
	GetThresholds() []CircuitBreakersThresholds
}

type CircuitBreakersThresholds interface {
}

// 如果某个服务实例频繁超时或者出错，则将该实例隔离，避免影响整个服务
type OutlierDetection interface {
	GetInterval() int64
}
