package svc

import (
	"github.com/xuexihuang/go-SaasSchedule/app/schedule/internal/config"
	"github.com/xuexihuang/go-SaasSchedule/app/schedule/internal/data"
)

type ServiceContext struct {
	Config        config.Config
	JobNodeRecord *data.JobNodeRecord
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		JobNodeRecord: data.NewJobNodeRecord(),
	}
}
