package biz

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/xuexihuang/go-SaasSchedule/app/schedule/internal/data/database"
	"github.com/xuexihuang/go-SaasSchedule/app/schedule/internal/svc"
	log15 "github.com/xuexihuang/new_log15"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

type NodeInter interface {
	RunSchedule(moduleId int64, jobId int64, chartUrl string, chartVersion string, domain string, imageTag string, tenantId string) error
	generateInitSql() string
	generateSetCommand(domain string, imageTag string, tenantId string) ([]string, error)
}

type JobNodeRecordRepo interface {
	Create(r *database.JobNodeRecord) (int64, error)
	UpdateStatus(id int64, status string) error
}

func NewJobNodeInter(moduleName string, chartVersion string, svcCtx *svc.ServiceContext) NodeInter {

	var ret NodeInter
	if moduleName == "user" {
		data := &UserNode{}
		data.moduleName = "user"
		data.nodeInter = data
		data.jobNodeRecordRepo = svcCtx.JobNodeRecord
		ret = data
	} else if moduleName == "admin" {
		data := &AdminNode{}
		data.moduleName = "admin"
		data.jobNodeRecordRepo = svcCtx.JobNodeRecord
		data.nodeInter = data
		ret = data
	} else {
		log15.Error("not find module impliment")
		return nil
	}
	return ret
}

type NodeBase struct {
	jobNodeRecordRepo JobNodeRecordRepo
	moduleName        string
	nodeInter         NodeInter
}

func (n *NodeBase) CreateRecord(r *database.JobNodeRecord) (int64, error) {
	return n.jobNodeRecordRepo.Create(r)
}
func (n *NodeBase) UpdateRecordStatus(id int64, status string) error {
	return n.jobNodeRecordRepo.UpdateStatus(id, status)
}
func (n *NodeBase) installChart(domain string, imageTag string, tenantId string, chartPath string, moduleId int64, jobId int64) (int64, error) {
	args, err := n.nodeInter.generateSetCommand(domain, imageTag, tenantId)
	if err != nil {
		log15.Error("generateSetCommand error")
		return 0, err
	}
	// 创建 helm install 命令
	cmdArgs := append([]string{"install", n.moduleName, n.moduleName + "/", "-f", n.moduleName + "/config.yaml", "--namespace", tenantId}, args...)
	cmd := exec.Command("helm", cmdArgs...)
	// 设置当前工作目录为 /data/gitCharts
	cmd.Dir = chartPath

	log15.Info("执行命令", "cmdArgs", cmdArgs)
	// 用于存储命令的输出
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	//插入数据库
	RecordId, err := n.CreateRecord(&database.JobNodeRecord{ModuleId: moduleId, JobId: jobId})
	if err != nil {
		return 0, err
	}
	// 运行命令并获取返回错误
	err = cmd.Run()

	log15.Info("cmd.Run", "err", err, "stderr", stderr.String(), "stdout", stdout.String())
	// 返回标准输出、错误输出和执行错误
	if err != nil {
		_ = n.UpdateRecordStatus(RecordId, "Error")
		return 0, err
	}
	return RecordId, nil
}
func (n *NodeBase) checkChartReleaseStatus(tenantId string) (string, string, error) {

	labelSelector := "app.kubernetes.io/instance=" + n.moduleName
	// 构建 kubectl 命令，不使用 JSON 输出
	cmd := exec.Command("kubectl", "get", "pods", "--namespace", tenantId, "-l", labelSelector)

	// 获取命令输出
	output, err := cmd.Output()
	if err != nil {
		return "", "", err
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
func (n *NodeBase) downChart(jobId int64, chartUrl string, chartVersion string) (string, error) {

	// 设置下载目录
	downloadDir := fmt.Sprintf("/data/gitCharts/%d", jobId)

	// 检查目标目录是否存在，如果不存在则创建
	if _, err := os.Stat(downloadDir); os.IsNotExist(err) {
		fmt.Println("Directory does not exist, creating it...")
		err := os.MkdirAll(downloadDir, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("failed to create directory %s: %v", downloadDir, err)
		}
	}

	// 构造 git clone 命令
	// 使用 --branch 或 -b 标记来指定分支或标签
	cmd := exec.Command("git", "clone", "--branch", chartVersion, chartUrl, downloadDir)

	// 执行命令并捕获输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to execute git clone: %v\nOutput: %s", err, output)
	}

	fmt.Println("Git clone completed successfully.")
	return downloadDir, nil
}
func (n *NodeBase) createK8sNamespace(tenantId string) error {

	// 构建 kubectl 命令，不使用 JSON 输出
	cmd := exec.Command("kubectl", "create", "ns", tenantId)
	// 获取命令输出
	_, err := cmd.Output()
	return err
}
func (n *NodeBase) RunSchedule(moduleId int64, jobId int64, chartUrl string, chartVersion string, domain string, imageTag string, tenantId string) error {

	log15.Info("RunSchedule", "moduledId", moduleId, "jobId", jobId, "chartUrl", chartUrl, "charversion", chartVersion, "domain", domain, "imageTag", imageTag, "tenantId", tenantId)
	chartPath, err := n.downChart(jobId, chartUrl, chartVersion)
	log15.Info("downChart", "chartPath", chartPath, "err", err)
	if err != nil {
		return err
	}
	err = n.createK8sNamespace(tenantId)
	log15.Info("createK8sNamespace", "err", err)
	if err != nil {
		return err
	}
	recordId, err := n.installChart(domain, imageTag, tenantId, chartPath, moduleId, jobId)
	log15.Info("installChart", "recordId", recordId, "err", err)
	if err != nil {
		return err
	}
	var successCount = 0
	for i := 0; i < 8; i++ {
		time.Sleep(2 * time.Second)
		podName, status, err := n.checkChartReleaseStatus(tenantId)
		log15.Info("checkChartReleaseStatus", "podName", podName, "status", status, "err", err)
		if err == nil && status == "Running" {
			successCount++
		} else {
			successCount = 0
		}
		if successCount >= 3 {
			break
		}
	}
	if successCount >= 3 {
		err = n.UpdateRecordStatus(recordId, "Running")
		log15.Info("UpdateRecordStatus", "err", err)
		return nil
	} else {
		err = n.UpdateRecordStatus(recordId, "Error")
		log15.Info("UpdateRecordStatus", "err", err)
		return errors.New("ChartReleaseStatus error")
	}

}
