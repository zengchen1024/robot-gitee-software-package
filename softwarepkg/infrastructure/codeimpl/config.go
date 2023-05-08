package codeimpl

type Config struct {
	ShellScript string      `json:"shell_script"`
	PkgSrcOrg   string      `json:"pkg_src_org"`
	Robot       RobotConfig `json:"robot"   required:"true"`
	CIRepo      CIRepo      `json:"ci_repo" required:"true"`
}

func (c *Config) SetDefault() {
	if c.ShellScript == "" {
		c.ShellScript = "/opt/app/code.sh"
	}

	if c.PkgSrcOrg == "" {
		c.PkgSrcOrg = "src-openeuler"
	}
}

type RobotConfig struct {
	Username string `json:"username" required:"true"`
	Token    string `json:"token"    required:"true"`
}

type CIRepo struct {
	Repo string `json:"repo" required:"true"`
	Link string `json:"link" required:"true"`
}
