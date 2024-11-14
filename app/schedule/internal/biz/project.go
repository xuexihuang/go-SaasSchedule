package biz

type ProjectDomain struct {
}
type Module struct {
	Id   int64
	Name string
	Sort int
}

func NewProjectDomain() *ProjectDomain {
	return &ProjectDomain{}
}
func (p *ProjectDomain) Get(id int64) ([]Module, string, error) {
	return []Module{{Id: 1, Name: "user", Sort: 1}, {Id: 2, Name: "admin", Sort: 2}, {Id: 3, Name: "game", Sort: 3}}, "git@github.com:xuexihuang/saas-chart.git", nil
}
