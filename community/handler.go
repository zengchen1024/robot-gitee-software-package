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
) *eventHandler {
	return &eventHandler{
		cli:     cli,
		cfg:     *cfg,
		repo:    repo,
		service: service,
	}
}

type eventHandler struct {
	cli     iClient
	cfg     Config
	repo    repository.SoftwarePkg
	service app.PackageService
}

func (impl *eventHandler) HandlePREvent(e *sdk.PullRequestEvent) error {
	if org, repo := e.GetOrgRepo(); !impl.cfg.isCommunity(org, repo) {
		return nil
	}

	if e.GetPullRequest().GetState() != sdk.StatusOpen {
		return impl.handlePRState(e)
	}

	if sdk.GetPullRequestAction(e) == sdk.PRActionUpdatedLabel {
		return impl.handleCILabel(e)
	}

	return nil
}
