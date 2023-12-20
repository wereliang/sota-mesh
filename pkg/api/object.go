package api

import "time"

// SotaObject
type SotaObject interface {
	UpdateTime() time.Time
	Config() interface{}
}
