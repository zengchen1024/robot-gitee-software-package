package pullrequestimpl

type Config struct {
	Robot       RobotConfig `json:"robot"`
	PR          PRConfig    `json:"pr"`
	Template    Template    `json:"template"`
	ShellScript string      `json:"shell_script"`
}

func (cfg *Config) SetDefault() {
	cfg.PR.setDefault()
	cfg.Template.setDefault()

	if cfg.ShellScript == "" {
		cfg.ShellScript = "/opt/app/repo.sh"
	}
}

type RobotConfig struct {
	Username string `json:"username" required:"true"`
	Token    string `json:"token"    required:"true"`
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

type Template struct {
	AppendSigInfo string `json:"append_sig_info"`
	NewRepoFile   string `json:"new_repo_file"`
}

func (t *Template) setDefault() {
	if t.AppendSigInfo == "" {
		t.AppendSigInfo = "/opt/app/template/append_sig_info.tpl"
	}

	if t.NewRepoFile == "" {
		t.NewRepoFile = "/opt/app/template/new_repo_file.tpl"
	}
}
