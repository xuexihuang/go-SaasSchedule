package logic

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/xuexihuang/go-SaasSchedule/app/schedule/internal/svc"
	"github.com/xuexihuang/go-SaasSchedule/app/schedule/internal/types"
	log15 "github.com/xuexihuang/new_log15"
	"os/exec"
	"regexp"
	"strings"
	"time"

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

// getPodStatus 通过命令行输出获取特定命名空间和标签的第一个 Pod 名称及状态
func getPodStatus(namespace, labelSelector string) (string, string, error) {
	// 构建 kubectl 命令，不使用 JSON 输出
	cmd := exec.Command("kubectl", "get", "pods", "--namespace", namespace, "-l", labelSelector)

	// 获取命令输出
	output, err := cmd.Output()
	if err != nil {
		return "", "", fmt.Errorf("failed to execute kubectl command: %v", err)
	}

	// 使用 bufio.Scanner 逐行解析输出
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	// 跳过第一行表头
	scanner.Scan()

	// 正则表达式匹配 Pod 名称和状态
	re := regexp.MustCompile(`^(\S+)\s+\S+\s+(\S+)`)

	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if len(matches) == 3 {
			podName := matches[1]
			podPhase := matches[2]
			return podName, podPhase, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", "", fmt.Errorf("failed to read command output: %v", err)
	}

	return "", "", fmt.Errorf("no pods found with the specified label selector")
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
	tStatus := ""
	for i := 0; i < 10; i++ {
		time.Sleep(2 * time.Second)
		labelSelector := "app.kubernetes.io/instance=" + req.Release

		podName, podStatus, err := getPodStatus(namespace, labelSelector)
		log15.Info("检测pod执行状态", "err", err, "podName", podName, "podStatus", podStatus)
		tStatus = podStatus
	}

	////////////////////////
	return &types.ShortenResp{Shorten: tStatus}, nil
}
