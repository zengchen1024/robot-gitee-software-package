package pullrequestimpl

type Config struct {
	Robot       RobotConfig `json:"robot"`
	PR          PRConfig    `json:"pr"`
	ShellScript string      `json:"shell_script"`
}

func (cfg *Config) SetDefault() {
	cfg.PR.setDefault()

	if cfg.ShellScript == "" {
		cfg.ShellScript = "./repo.sh"
	}
}

type RobotConfig struct {
	Username string `json:"username" required:"true"`
	Password string `json:"password" required:"true"`
	Email    string `json:"email"    required:"true"`
}

type PRConfig struct {
	NewRepoBranch NewRepoBranch `json:"new_repo_branch"`
	Org           string        `json:"org"`
	Repo          string        `json:"repo"`
	PRName        string        `json:"pr_name"`
}

func (cfg *PRConfig) setDefault() {
	if cfg.Org == "" {
		cfg.Org = "openeuler"
	}

	if cfg.Repo == "" {
		cfg.Repo = "community"
	}

	if cfg.PRName == "" {
		cfg.PRName = ",新增软件包申请"
	}

	cfg.NewRepoBranch.setDefault()
}

type NewRepoBranch struct {
	Name        string `json:"name"`
	ProtectType string `json:"protect_type"`
	PublicType  string `json:"public_type"`
}

func (cfg *NewRepoBranch) setDefault() {
	if cfg.Name == "" {
		cfg.Name = "master"
	}

	if cfg.ProtectType == "" {
		cfg.ProtectType = "protected"
	}

	if cfg.PublicType == "" {
		cfg.PublicType = "public"
	}
}
