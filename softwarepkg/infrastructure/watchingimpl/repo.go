package watchingimpl

import (
	"time"

	sdk "github.com/opensourceways/go-gitee/gitee"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/app"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain/repository"
)

func NewWatchingImpl(cfg Config, cli iClient) *WatchingImpl {
	return &WatchingImpl{
		cfg: cfg,
		cli: cli,
	}
}

type iClient interface {
	GetRepo(org, repo string) (sdk.Project, error)
}

type WatchingImpl struct {
	cfg       Config
	cli       iClient
	repo      repository.PullRequest
	prService app.PullRequestService
}

func (impl *WatchingImpl) Run() {
	interval := impl.cfg.IntervalDuration()

	for {
		prs, err := impl.repo.FindAll(true)
		if err != nil {
			logrus.Errorf("find all storage pr failed, err: %s", err.Error())
		}

		for _, pr := range prs {
			if !pr.IsMerged() {
				continue
			}

			v, err := impl.cli.GetRepo(impl.cfg.Org, pr.Pkg.Name)
			if err != nil {
				continue
			}

			if err = impl.prService.HandleRepoCreated(&pr, v.Url); err != nil {
				logrus.Errorf("handle repo created err: %s", err.Error())
			}
		}

		time.Sleep(interval)
	}
}
