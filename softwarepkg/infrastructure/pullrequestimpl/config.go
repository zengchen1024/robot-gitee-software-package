package pullrequestimpl

type Config struct {
	CommunityRobot CommunityRobotConfig `json:"community_robot"`
	Robot          RobotConfig          `json:"robot"`
	Template       Template             `json:"template"`
	ShellScript    string               `json:"shell_script"`
	SoftwarePkg    SoftwarePkg          `json:"software_pkg"`
}

func (cfg *Config) SetDefault() {
	cfg.Robot.setDefault()
	cfg.CommunityRobot.setDefault()
	cfg.Template.setDefault()

	if cfg.ShellScript == "" {
		cfg.ShellScript = "/opt/app/repo.sh"
	}
}

type RobotConfig struct {
	Username      string        `json:"username" required:"true"`
	Token         string        `json:"token"    required:"true"`
	Email         string        `json:"email"    required:"true"`
	Repo          string        `json:"repo"     required:"true"`
	NewRepoBranch NewRepoBranch `json:"new_repo_branch"`
}

func (cfg *RobotConfig) setDefault() {
	cfg.NewRepoBranch.setDefault()
}

type CommunityRobotConfig struct {
	Token string `json:"token" required:"true"`
	Org   string `json:"org"`
	Repo  string `json:"repo"`
}

func (cfg *CommunityRobotConfig) setDefault() {
	if cfg.Org == "" {
		cfg.Org = "openeuler"
	}

	if cfg.Repo == "" {
		cfg.Repo = "community"
	}
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
	PRBodyTpl   string `json:"pr_body_tpl"`
	SigInfoTpl  string `json:"sig_info_tpl"`
	RepoYamlTpl string `json:"repo_yaml_tpl"`
}

func (t *Template) setDefault() {
	if t.PRBodyTpl == "" {
		t.PRBodyTpl = "/opt/app/template/pr_body.tpl"
	}

	if t.SigInfoTpl == "" {
		t.SigInfoTpl = "/opt/app/template/sig_info.tpl"
	}

	if t.RepoYamlTpl == "" {
		t.RepoYamlTpl = "/opt/app/template/repo_yaml.tpl"
	}
}

type SoftwarePkg struct {
	Endpoint string `json:"endpoint" required:"true"`
}
