package biz

type Paths struct {
	Path     string `yaml:"path"`
	PathType string `yaml:"pathType"`
}
type Hosts struct {
	Host  string   `yaml:"host"`
	Paths []*Paths `yaml:"paths"`
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
		Host  string `yaml:"Host"`
		Port  int    `yaml:"Port"`
		Mysql struct {
			DataSource string `yaml:"Port"`
		}
	}
	SqlConfig struct {
		Sql string `yaml:"sql"`
	}
}
