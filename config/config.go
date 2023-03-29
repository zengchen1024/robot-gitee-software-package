package config

import (
	"github.com/opensourceways/server-common-lib/utils"

	"github.com/opensourceways/robot-gitee-software-package/kafka"
	"github.com/opensourceways/robot-gitee-software-package/message-server"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/infrastructure/codeimpl"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/infrastructure/emailimpl"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/infrastructure/postgresql"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/infrastructure/pullrequestimpl"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/infrastructure/repositoryimpl"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/infrastructure/watchingimpl"
	localutils "github.com/opensourceways/robot-gitee-software-package/utils"
)

func LoadConfig(path string) (*Config, error) {
	cfg := new(Config)
	if err := utils.LoadFromYaml(path, cfg); err != nil {
		return nil, err
	}

	cfg.SetDefault()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

type configValidate interface {
	Validate() error
}

type configSetDefault interface {
	SetDefault()
}

type PostgresqlConfig struct {
	DB postgresql.Config `json:"db" required:"true"`

	repositoryimpl.Config
}

type Config struct {
	MQ            kafka.Config           `json:"mq"`
	MessageServer messageserver.Config   `json:"message_server"`
	Email         emailimpl.Config       `json:"email"`
	Watch         watchingimpl.Config    `json:"watch"`
	Postgresql    PostgresqlConfig       `json:"postgresql"`
	PullRequest   pullrequestimpl.Config `json:"pull_request"`
	Code          codeimpl.Config        `json:"code"`
	Encryption    localutils.Config      `json:"encryption"`
}

func (cfg *Config) configItems() []interface{} {
	return []interface{}{
		&cfg.MQ,
		&cfg.MessageServer,
		&cfg.Email,
		&cfg.Watch,
		&cfg.Postgresql.DB,
		&cfg.Postgresql.Config,
		&cfg.PullRequest,
		&cfg.Code,
		&cfg.Encryption,
	}
}

func (cfg *Config) SetDefault() {
	items := cfg.configItems()
	for _, i := range items {
		if f, ok := i.(configSetDefault); ok {
			f.SetDefault()
		}
	}
}

func (cfg *Config) Validate() error {
	if _, err := utils.BuildRequestBody(cfg, ""); err != nil {
		return err
	}

	items := cfg.configItems()
	for _, i := range items {
		if f, ok := i.(configValidate); ok {
			if err := f.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}
