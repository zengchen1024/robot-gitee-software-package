package main

import (
	"errors"

	sdk "github.com/opensourceways/go-gitee/gitee"
	"github.com/opensourceways/server-common-lib/config"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/robot-gitee-lib/framework"
)

// TODO: set botName
const botName = "software-package"

type iClient interface {
}

func newRobot(cli iClient) *robot {
	return &robot{cli: cli}
}

type robot struct {
	cli iClient
}

func (bot *robot) NewConfig() config.Config {
	return &configuration{}
}

func (bot *robot) getConfig(cfg config.Config) (*configuration, error) {
	if c, ok := cfg.(*configuration); ok {
		return c, nil
	}
	return nil, errors.New("can't convert to configuration")
}

func (bot *robot) RegisterEventHandler(f framework.HandlerRegister) {
	f.RegisterPullRequestHandler(bot.handlePREvent)
}

func (bot *robot) handlePREvent(e *sdk.PullRequestEvent, c config.Config, log *logrus.Entry) error {
	// TODO: if it doesn't needd to hand PR event, delete this function.
	return nil
}
