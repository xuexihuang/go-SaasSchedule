package data

import (
	"errors"
	"github.com/xuexihuang/go-SaasSchedule/app/schedule/internal/data/database"
	"github.com/xuexihuang/go-SaasSchedule/pkg/data"
	"time"
)

type JobNodeRecord struct {
}

func NewJobNodeRecord() *JobNodeRecord {
	return &JobNodeRecord{}
}
func (j *JobNodeRecord) Create(r *database.JobNodeRecord) (int64, error) {
	r.StartedAt = time.Now().Unix()
	r.Status = 1
	_, err := data.MysqlEngine.Table(&database.JobNodeRecord{}).Insert(r)
	if err != nil {
		return 0, err
	}
	return r.Id, nil
}
func (j *JobNodeRecord) UpdateStatus(id int64, status string) error {
	m := &database.JobNodeRecord{Id: id}
	if status == "Running" {
		m.Status = 2
	} else {
		m.Status = 3
	}
	_, err := data.MysqlEngine.ID(id).Cols("status").Update(m)
	if err != nil {
		return errors.New("update error")
	}
	return nil
}
