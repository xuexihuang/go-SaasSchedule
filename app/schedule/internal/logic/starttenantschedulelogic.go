package logic

import (
	"context"
	"errors"
	"github.com/xuexihuang/go-SaasSchedule/app/schedule/internal/biz"
	log15 "github.com/xuexihuang/new_log15"

	"github.com/xuexihuang/go-SaasSchedule/app/schedule/internal/svc"
	"github.com/xuexihuang/go-SaasSchedule/app/schedule/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type StartTenantScheduleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStartTenantScheduleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StartTenantScheduleLogic {
	return &StartTenantScheduleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StartTenantScheduleLogic) StartTenantSchedule(req *types.StartTenantScheduleReq) (resp *types.StartTenantScheduleResp, err error) {

	modules, chartUrl, err := biz.NewProjectDomain().Get(req.ProjectId)
	if err != nil {
		return nil, errors.New("获取项目模块详情错误")
	}
	jobId, err := biz.NewJobDomain().CreateScheduleJob(req.ProjectId, req.Domain, req.ChartVersion, req.ImageTag, req.TenantId)
	if err != nil {
		return nil, errors.New("CreateScheduleJob error")
	}
	go func() { //异步进行自动化部署
		for _, v := range modules {
			nodeInter := biz.NewJobNodeInter(v.Name, req.ChartVersion, l.svcCtx)
			err := nodeInter.RunSchedule(jobId, chartUrl, req.ChartVersion, req.Domain, req.ImageTag, req.TenantId)
			if err != nil {
				log15.Error("RunSchedule error", "err", err)
				break
			}
		}
	}()

	return &types.StartTenantScheduleResp{JobId: jobId}, nil
}
