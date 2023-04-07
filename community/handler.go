package community

import (
	sdk "github.com/opensourceways/go-gitee/gitee"
	"github.com/opensourceways/robot-gitee-lib/client"

	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/app"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain/repository"
)

type EventHandler interface {
	HandlePREvent(e *sdk.PullRequestEvent) error
}

type iClient interface {
	GetRepo(org, repo string) (sdk.Project, error)
}

func NewEventHandler(
	cfg *Config,
	repo repository.SoftwarePkg,
	service app.PackageService,
) (*eventHandler, error) {
	cli := client.NewClient(func() []byte {
		return []byte(cfg.RobotToken)
	})

	u, err := cli.GetBot()
	if err != nil {
		return nil, err
	}

	return &eventHandler{
		cli: robotImpl{
			iClient:   cli,
			robotName: u.Login,
		},
		cfg:     *cfg,
		repo:    repo,
		service: service,
	}, nil
}

type robotImpl struct {
	iClient
	robotName string
}

type eventHandler struct {
	cli     robotImpl
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
