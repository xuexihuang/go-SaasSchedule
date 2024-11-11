package logic

import (
	"bytes"
	"context"
	"fmt"
	log15 "github.com/xuexihuang/new_log15"
	"os/exec"

	"github.com/xuexihuang/go-SaasSchedule/app/schedule/internal/svc"
	"github.com/xuexihuang/go-SaasSchedule/app/schedule/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShortenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShortenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShortenLogic {
	return &ShortenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
func installChart(releaseName, chartPath string, namespace string, args ...string) (string, error) {
	// 创建 helm install 命令
	cmdArgs := append([]string{"install", releaseName, chartPath, "--namespace", namespace}, args...)
	cmd := exec.Command("helm", cmdArgs...)

	// 用于存储命令的输出
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// 运行命令并获取返回错误
	err := cmd.Run()

	// 返回标准输出、错误输出和执行错误
	if err != nil {
		return stderr.String(), fmt.Errorf("command failed with error: %v", err)
	}
	return stdout.String(), nil
}
func (l *ShortenLogic) Shorten(req *types.ShortenReq) (resp *types.ShortenResp, err error) {

	log15.Info("进入shorten调用", "req", req)
	// 示例：安装 chart
	releaseName := req.Release
	chartPath := "."
	namespace := req.Name

	// 执行 helm install
	output, err := installChart(releaseName, chartPath, namespace, "-f", "config.yaml")
	if err != nil {
		log15.Error("执行helm输出错误", "err", err, "out", output)
	} else {
		log15.Info("执行helm输出正常", "out", output)
	}
	return
}
