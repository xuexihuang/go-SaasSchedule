package biz

import "time"

type JobDomain struct {
}

func NewJobDomain() *JobDomain {
	return &JobDomain{}
}
func (j *JobDomain) CreateScheduleJob(projectId int64, domain string, chartVersion string, imageTag string, tenantId string) (int64, error) {
	return time.Now().Unix(), nil
}
func (j *JobDomain) UpdateScheduleJobStatus(id int64, status int) error {

	return nil
}
