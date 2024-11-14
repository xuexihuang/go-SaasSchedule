package biz

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/xuexihuang/go-SaasSchedule/app/schedule/internal/data/database"
	"github.com/xuexihuang/go-SaasSchedule/app/schedule/internal/svc"
	log15 "github.com/xuexihuang/new_log15"
	"gopkg.in/yaml.v2"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

type NodeInter interface {
	RunSchedule(moduleId int64, jobId int64, chartUrl string, chartVersion string, domain string, imageTag string, tenantId string) error
	generateInitSql() string
	generateSetCommand(domain string, imageTag string, tenantId string) (interface{}, error)
}

type JobNodeRecordRepo interface {
	Create(r *database.JobNodeRecord) (int64, error)
	UpdateStatus(id int64, status string) error
}

func NewJobNodeInter(moduleName string, chartVersion string, svcCtx *svc.ServiceContext) NodeInter {

	var ret NodeInter
	if moduleName == "user" {
		data := &UserNode{NodeBase: &NodeBase{}}
		data.moduleName = "user"
		data.nodeInter = data
		data.jobNodeRecordRepo = svcCtx.JobNodeRecord
		ret = data
	} else if moduleName == "admin" {
		data := &AdminNode{NodeBase: &NodeBase{}}
		data.moduleName = "admin"
		data.jobNodeRecordRepo = svcCtx.JobNodeRecord
		data.nodeInter = data
		ret = data
	} else if moduleName == "game" {
		data := &GameNode{NodeBase: &NodeBase{}}
		data.moduleName = "game"
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
	c, err := n.nodeInter.generateSetCommand(domain, imageTag, tenantId)
	if err != nil {
		log15.Error("generateSetCommand error")
		return 0, err
	}
	// 序列化为 YAML
	data, err := yaml.Marshal(c)
	if err != nil {
		log15.Error("yaml.Marshal error", "err", err)
		return 0, err
	}
	// 打印 YAML 输出
	fmt.Println(string(data))
	// 将 YAML 数据写入到文件
	err = os.WriteFile(chartPath+"/"+n.moduleName+"/config.yaml", data, 0644)
	if err != nil {
		log15.Error("failed to write to file", "err", err)
		return 0, err
	}

	// 创建 helm install 命令
	cmdArgs := []string{"install", n.moduleName, n.moduleName + "/", "-f", n.moduleName + "/config.yaml", "--namespace", tenantId}
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
	_, err := os.Stat(downloadDir)
	if os.IsNotExist(err) {
		fmt.Println("Directory does not exist, creating it...")
		err := os.MkdirAll(downloadDir, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("failed to create directory %s: %v", downloadDir, err)
		}
	} else if err != nil {

		return "", fmt.Errorf("检测目录 %s 时出错: %v\n", downloadDir, err)
	} else {
		fmt.Printf("目录 %s 已存在\n", downloadDir)
		return downloadDir, nil
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

	// 先执行 `kubectl get ns` 命令
	cmd := exec.Command("kubectl", "get", "ns")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("获取命名空间列表失败: %v", err)
	}

	// 使用正则表达式检查是否存在相同的 tenantId
	regex := regexp.MustCompile(fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta(tenantId)))
	if regex.Match(out.Bytes()) {
		fmt.Printf("命名空间 %s 已存在，跳过创建。\n", tenantId)
		return nil
	}

	// 如果不存在，则创建命名空间
	cmd = exec.Command("kubectl", "create", "ns", tenantId)
	_, err = cmd.Output()
	if err != nil {
		return fmt.Errorf("创建命名空间 %s 失败: %v", tenantId, err)
	}

	fmt.Printf("命名空间 %s 创建成功。\n", tenantId)
	return nil
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
	for i := 0; i < 18; i++ {
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
		log15.Info("恭喜你成功部署了saas系统！\n")
		return nil
	} else {
		err = n.UpdateRecordStatus(recordId, "Error")
		log15.Info("UpdateRecordStatus", "err", err)
		return errors.New("ChartReleaseStatus error")
	}

}
