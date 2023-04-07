package watch

import (
	"time"

	sdk "github.com/opensourceways/go-gitee/gitee"
	"github.com/opensourceways/robot-gitee-lib/client"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/app"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain/repository"
)

func NewWatchingImpl(
	cfg Config,
	repo repository.SoftwarePkg,
	service app.PackageService,
) *WatchingImpl {
	cli := client.NewClient(func() []byte {
		return []byte(cfg.RobotToken)
	})

	return &WatchingImpl{
		cfg:     cfg,
		cli:     cli,
		repo:    repo,
		service: service,
		stop:    make(chan struct{}),
		stopped: make(chan struct{}),
	}
}

type iClient interface {
	GetRepo(org, repo string) (sdk.Project, error)
}

type WatchingImpl struct {
	cfg     Config
	cli     iClient
	repo    repository.SoftwarePkg
	service app.PackageService
	stop    chan struct{}
	stopped chan struct{}
}

func (impl *WatchingImpl) Start() {
	go impl.watch()
}

func (impl *WatchingImpl) Stop() {
	close(impl.stop)

	<-impl.stopped
}

func (impl *WatchingImpl) watch() {
	interval := impl.cfg.IntervalDuration()

	checkStop := func() bool {
		select {
		case <-impl.stop:
			return true
		default:
			return false
		}
	}

	for {
		prs, err := impl.repo.FindAll()
		if err != nil {
			logrus.Errorf("find all storage pr failed, err: %s", err.Error())
		}

		for _, pr := range prs {
			impl.handle(pr)

			if checkStop() {
				close(impl.stopped)

				return
			}
		}

		time.Sleep(interval)
	}
}

func (impl *WatchingImpl) handle(pkg domain.SoftwarePkg) {
	switch pkg.Status {
	case domain.PkgStatusPRMerged:
		v, err := impl.cli.GetRepo(impl.cfg.PkgOrg, pkg.Name)
		if err != nil {
			return
		}

		if err = impl.service.HandleRepoCreated(&pkg, v.HtmlUrl); err != nil {
			logrus.Errorf("handle repo created err: %s", err.Error())
		}

	case domain.PkgStatusRepoCreated:
		if err := impl.service.HandlePushCode(&pkg); err != nil {
			logrus.Errorf("handle push code err: %s", err.Error())
		}
	}
}
