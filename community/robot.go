package community

import (
	sdk "github.com/opensourceways/go-gitee/gitee"

	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/app"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain/repository"
)

type EventHandler interface {
	HandlePREvent(e *sdk.PullRequestEvent) error
}

type iClient interface {
	GetBot() (sdk.User, error)
	GetRepo(org, repo string) (sdk.Project, error)
}

func NewEventHandler(
	cli iClient,
	cfg *Config,
	repo repository.SoftwarePkg,
	service app.PackageService,
) *robot {
	return &robot{
		cli:     cli,
		cfg:     *cfg,
		repo:    repo,
		service: service,
	}
}

type robot struct {
	cli     iClient
	cfg     Config
	repo    repository.SoftwarePkg
	service app.PackageService
}

func (bot *robot) HandlePREvent(e *sdk.PullRequestEvent) error {
	org, repo := e.GetOrgRepo()
	if bot.cfg.isCommunity(org, repo) {
		return nil
	}

	prState := e.GetPullRequest().GetState()

	if prState != sdk.StatusOpen {
		return bot.handlePRState(e)
	}

	if sdk.GetPullRequestAction(e) != sdk.PRActionUpdatedLabel {
		return nil
	}

	return bot.handleCILabel(e)
}
