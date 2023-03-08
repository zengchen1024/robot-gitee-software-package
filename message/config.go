package message

import (
	"errors"

	"github.com/opensourceways/server-common-lib/utils"
)

type config struct {
	KafkaAddress string      `json:"kafka_address"   required:"true"`
	Topics       Topics      `json:"topics"`
	Robot        RobotConfig `json:"robot"`
	PR           PRConfig    `json:"pr"`
}

func loadConfig(path string) (*config, error) {
	cfg := new(config)
	if err := utils.LoadFromYaml(path, cfg); err != nil {
		return nil, err
	}

	cfg.SetDefault()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *config) Validate() error {
	if c.KafkaAddress == "" {
		return errors.New("missing kafka_address")
	}

	if c.Topics.NewPkg == "" {
		return errors.New("missing new pkg topic")
	}

	if c.Topics.CIPassed == "" {
		return errors.New("missing ci passed topic")
	}

	if c.Robot.Username == "" {
		return errors.New("missing robot username")
	}

	if c.Robot.Password == "" {
		return errors.New("missing robot password")
	}

	if c.Robot.Email == "" {
		return errors.New("missing robot email")
	}

	return nil
}

func (c *config) SetDefault() {
	if c.PR.NewRepoBranch.Name == "" {
		c.PR.NewRepoBranch.Name = "master"
	}

	if c.PR.NewRepoBranch.ProtectType == "" {
		c.PR.NewRepoBranch.ProtectType = "protected"
	}

	if c.PR.NewRepoBranch.PublicType == "" {
		c.PR.NewRepoBranch.PublicType = "public"
	}

	if c.PR.Org == "" {
		c.PR.Org = "openeuler"
	}

	if c.PR.Repo == "" {
		c.PR.Repo = "community"
	}

	if c.PR.BranchName == "" {
		c.PR.BranchName = "software_pkg_%s"
	}

	if c.PR.PRName == "" {
		c.PR.PRName = "software_pkg_%s,新增软件包申请"
	}

	if c.PR.ModifyFiles.SigInfo == "" {
		c.PR.ModifyFiles.SigInfo = "community/sig/%s/sig-info.yaml"
	}

	if c.PR.ModifyFiles.NewRepo == "" {
		c.PR.ModifyFiles.NewRepo = "community/sig/%s/src-openeuler/%s/%s.yaml"
	}
}

type RobotConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Topics struct {
	NewPkg   string `json:"new_pkg"`
	CIPassed string `json:"ci_passed"`
}

type PRConfig struct {
	ModifyFiles   ModifyFiles   `json:"modify_files"`
	NewRepoBranch NewRepoBranch `json:"new_repo_branch"`
	Org           string        `json:"org"`
	Repo          string        `json:"repo"`
	BranchName    string        `json:"branch_name"`
	PRName        string        `json:"pr_name"`
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
