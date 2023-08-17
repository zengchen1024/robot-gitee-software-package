package watch

import "time"

type Config struct {
	RobotToken     string `json:"robot_token"      required:"true"`
	PkgOrg         string `json:"pkg_org"          required:"true"`
	CommunityOrg   string `json:"community_org"    required:"true"`
	CommunityRepo  string `json:"community_repo"   required:"true"`
	CISuccessLabel string `json:"ci_success_label" required:"true"`
	CIFailureLabel string `json:"ci_failure_label" required:"true"`
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
