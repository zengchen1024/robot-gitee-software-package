package config

import (
	kafka "github.com/opensourceways/kafka-lib/agent"
	"github.com/opensourceways/server-common-lib/utils"

	"github.com/opensourceways/robot-gitee-software-package/message-server"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/infrastructure/codeimpl"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/infrastructure/emailimpl"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/infrastructure/postgresql"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/infrastructure/pullrequestimpl"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/infrastructure/repositoryimpl"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/watch"
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
	Code          codeimpl.Config        `json:"code"`
	Kafka         kafka.Config           `json:"kafka"`
	Email         emailimpl.Config       `json:"email"`
	Watch         watch.Config           `json:"watch"`
	Postgresql    PostgresqlConfig       `json:"postgresql"`
	Encryption    localutils.Config      `json:"encryption"`
	PullRequest   pullrequestimpl.Config `json:"pull_request"`
	MessageServer messageserver.Config   `json:"message_server"`
}

func (cfg *Config) configItems() []interface{} {
	return []interface{}{
		&cfg.Code,
		&cfg.Kafka,
		&cfg.Email,
		&cfg.Watch,
		&cfg.Postgresql.DB,
		&cfg.Postgresql.Config,
		&cfg.Encryption,
		&cfg.PullRequest,
		&cfg.MessageServer,
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
