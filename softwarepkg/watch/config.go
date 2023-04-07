package watch

import "time"

type Config struct {
	RobotToken string `json:"robot_token" required:"true"`
	Org        string `json:"org"`
	// unit second
	Interval int `json:"interval"`
}

func (cfg *Config) SetDefault() {
	if cfg.Org == "" {
		cfg.Org = "src-openeuler"
	}

	if cfg.Interval <= 0 {
		cfg.Interval = 10
	}
}

func (cfg *Config) IntervalDuration() time.Duration {
	return time.Second * time.Duration(cfg.Interval)
}
