package svc

import (
	"github.com/xuexihuang/go-SaasSchedule/app/schedule/internal/config"
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
