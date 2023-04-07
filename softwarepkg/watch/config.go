package watch

import "time"

type Config struct {
	RobotToken string `json:"robot_token" required:"true"`
	PkgOrg     string `json:"pkg_org"`
	// unit second
	Interval int `json:"interval"`
}

func (cfg *Config) SetDefault() {
	if cfg.PkgOrg == "" {
		cfg.PkgOrg = "src-openeuler"
	}

	if cfg.Interval <= 0 {
		cfg.Interval = 10
	}
}

func (cfg *Config) IntervalDuration() time.Duration {
	return time.Second * time.Duration(cfg.Interval)
}
