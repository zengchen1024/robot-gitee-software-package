package main

import (
	"fmt"

	sdk "github.com/opensourceways/go-gitee/gitee"
	"github.com/opensourceways/robot-gitee-lib/framework"
	"github.com/opensourceways/server-common-lib/config"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/app"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain/repository"
)

// TODO: set botName
const botName = "software-package"

type iClient interface {
	GetBot() (sdk.User, error)
	GetRepo(org, repo string) (sdk.Project, error)
}

func newRobot(
	cli iClient, prService app.PullRequestService,
	r repository.PullRequest, org string,
) *robot {
	return &robot{
		cli:       cli,
		prService: prService,
		repo:      r,
		PkgSrcOrg: org,
	}
}

type robot struct {
	cli       iClient
	prService app.PullRequestService
	repo      repository.PullRequest
	PkgSrcOrg string
}

func (bot *robot) NewConfig() config.Config {
	return &configuration{}
}

func (bot *robot) getConfig(cfg config.Config, org, repo string) (*botConfig, error) {
	c, ok := cfg.(*configuration)
	if !ok {
		return nil, fmt.Errorf("can't convert to configuration")
	}

	if bc := c.configFor(org, repo); bc != nil {
		return bc, nil
	}

	return nil, fmt.Errorf("no config for this repo:%s/%s", org, repo)
}

func (bot *robot) RegisterEventHandler(f framework.HandlerRegister) {
	f.RegisterPullRequestHandler(bot.handlePREvent)
}

func (bot *robot) handlePREvent(e *sdk.PullRequestEvent, c config.Config, log *logrus.Entry) error {
	org, repo := e.GetOrgRepo()
	cfg, err := bot.getConfig(c, org, repo)
	if err != nil {
		return err
	}

	prState := e.GetPullRequest().GetState()

	if prState != sdk.StatusOpen {
		return bot.handlePRState(e)
	}

	if sdk.GetPullRequestAction(e) != sdk.PRActionUpdatedLabel {
		return nil
	}

	return bot.handleCILabel(e, cfg)
}
