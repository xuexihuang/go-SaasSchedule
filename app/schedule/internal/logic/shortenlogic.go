package logic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	log15 "github.com/xuexihuang/new_log15"
	"os/exec"
	"time"

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

	log15.Info("执行命令", "cmdArgs", cmdArgs)
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

// PodStatus 用于解析命令返回的 Pod 状态
type PodStatus struct {
	Items []struct {
		Metadata struct {
			Name string `json:"name"`
		} `json:"metadata"`
		Status struct {
			Phase string `json:"phase"`
		} `json:"status"`
	} `json:"items"`
}

// getPodStatus 获取特定命名空间和标签的第一个 Pod 名称及状态
func getPodStatus(namespace, labelSelector string) (string, string, error) {
	// 构建 kubectl 命令
	cmd := exec.Command("kubectl", "get", "pods", "--namespace", namespace, "-l", labelSelector, "-o", "json")

	// 获取命令输出
	output, err := cmd.Output()
	if err != nil {
		return "", "", fmt.Errorf("failed to execute kubectl command: %v", err)
	}

	// 解析输出 JSON
	var podStatus PodStatus
	if err := json.Unmarshal(output, &podStatus); err != nil {
		return "", "", fmt.Errorf("failed to parse JSON output: %v", err)
	}

	// 检查是否有 Pod 返回
	if len(podStatus.Items) == 0 {
		return "", "", fmt.Errorf("no pods found with the specified label selector")
	}

	// 获取第一个 Pod 的名称和状态
	podName := podStatus.Items[0].Metadata.Name
	podPhase := podStatus.Items[0].Status.Phase
	return podName, podPhase, nil
}

func (l *ShortenLogic) Shorten(req *types.ShortenReq) (resp *types.ShortenResp, err error) {

	log15.Info("进入shorten调用", "req", req)
	// 示例：安装 chart
	releaseName := req.Release
	chartPath := "/data/github/helm-charts/testweb/"
	namespace := req.Name

	// 执行 helm install
	output, err := installChart(releaseName, chartPath, namespace, "-f", "/data/github/helm-charts/testweb/config.yaml")
	if err != nil {
		log15.Error("执行helm输出错误", "err", err, "out", output)
	} else {
		log15.Info("执行helm输出正常", "out", output)
	}
	////////////////////////////
	time.Sleep(2 * time.Second)
	labelSelector := "app.kubernetes.io/instance=" + req.Release

	podName, podStatus, err := getPodStatus(namespace, labelSelector)
	log15.Info("检测pod执行状态", "err", err, "podName", podName, "podStatus", podStatus)
	////////////////////////
	return
}
