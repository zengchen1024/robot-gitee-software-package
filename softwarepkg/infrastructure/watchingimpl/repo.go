package watchingimpl

import (
	"context"
	"time"

	sdk "github.com/opensourceways/go-gitee/gitee"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/app"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain/repository"
)

func NewWatchingImpl(
	cfg Config, cli iClient,
	repo repository.PullRequest, prService app.PullRequestService,
) *WatchingImpl {
	return &WatchingImpl{
		cfg:       cfg,
		cli:       cli,
		repo:      repo,
		prService: prService,
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

func (impl *WatchingImpl) Start(ctx context.Context, stop chan struct{}) {
	interval := impl.cfg.IntervalDuration()

	checkStop := func() bool {
		select {
		case <-ctx.Done():
			return true
		default:
			return false
		}
	}

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

			if err = impl.prService.HandleRepoCreated(&pr, v.HtmlUrl); err != nil {
				logrus.Errorf("handle repo created err: %s", err.Error())
			} else {
				if err = impl.prService.HandlePushCode(&pr); err != nil {
					logrus.Errorf("handle push code err: %s", err.Error())
				}
			}

			if checkStop() {
				close(stop)

				return
			}
		}

		time.Sleep(interval)
	}
}
