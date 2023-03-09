package message

type Config struct {
	Robot RobotConfig `json:"robot"`
	PR    PRConfig    `json:"pr"`
}

func (cfg *Config) SetDefault() {
	cfg.PR.setDefault()

	/*
		if pr.BranchName == "" {
			pr.BranchName = "software_pkg_%s"
		}

		if pr.PRName == "" {
			pr.PRName = "software_pkg_%s,新增软件包申请"
		}

		if pr.ModifyFiles.SigInfo == "" {
			pr.ModifyFiles.SigInfo = "community/sig/%s/sig-info.yaml"
		}

		if pr.ModifyFiles.NewRepo == "" {
			pr.ModifyFiles.NewRepo = "community/sig/%s/src-openeuler/%s/%s.yaml"
		}
	*/
}

type RobotConfig struct {
	Username string `json:"username" required:"true"`
	Password string `json:"password" required:"true"`
	Email    string `json:"email"    required:"true"`
}

type PRConfig struct {
	ModifyFiles   ModifyFiles   `json:"modify_files"`
	NewRepoBranch NewRepoBranch `json:"new_repo_branch"`
	Org           string        `json:"org"`
	Repo          string        `json:"repo"`
	PRName        string        `json:"pr_name"`
	BranchName    string        `json:"branch_name"`
}

func (cfg *PRConfig) setDefault() {
	if cfg.Org == "" {
		cfg.Org = "openeuler"
	}

	if cfg.Repo == "" {
		cfg.Repo = "community"
	}

	cfg.NewRepoBranch.setDefault()
}

type ModifyFiles struct {
	SigInfo string `json:"sig_info"`
	NewRepo string `json:"new_repo"`
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
