package biz

type Paths struct {
	Path     string `yaml:"path"`
	PathType string `yaml:"pathType"`
}
type Hosts struct {
	Host  string   `yaml:"host"`
	Paths []*Paths `yaml:"paths"`
}
type SqlConfig struct {
	Sql string `yaml:"sql"`
}

// 定义结构体
type UserConfig struct {
	Image struct {
		Tag string `yaml:"tag"`
	}
	Ingress struct {
		Enabled bool     `yaml:"enabled"`
		Hosts   []*Hosts `yaml:"hosts"`
	}
	Config struct {
		Host  string `yaml:"host"`
		Port  int    `yaml:"port"`
		Mysql struct {
			DataSource string `yaml:"dataSource"`
		}
	}
	SqlConfig SqlConfig `yaml:"sqlConfig"`
}
