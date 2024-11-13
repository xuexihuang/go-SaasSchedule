package database

type (
	JobNodeRecord struct {
		Id        int64  `xorm:"pk autoincr 'id'"`
		ModuleId  int64  `xorm:"module_id"`
		JobId     int64  `xorm:"job_id"`
		Status    int64  `xorm:"job_id"`
		ErrMsg    string `xorm:"err_msg"`
		StartedAt int64  `xorm:"started_at"`
		EndedAt   int64  `xorm:"ended_at"`
	}
)

func (JobNodeRecord) TableName() string {
	return "job_node_record"
}
