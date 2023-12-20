package config

type LoadBalancerType int32

const (
	Round_Robin     LoadBalancerType = 0
	LB_Original_Dst LoadBalancerType = 10
	LB_Logical_DNS  LoadBalancerType = 11
)
