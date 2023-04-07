package pullrequestimpl

import (
	sdk "github.com/opensourceways/go-gitee/gitee"
	"github.com/opensourceways/robot-gitee-lib/client"

	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain"
)

func NewPullRequestImpl(cfg *Config) (*pullRequestImpl, error) {
	cli := client.NewClient(func() []byte {
		return []byte(cfg.Robot.Token)
	})

	robot := client.NewClient(func() []byte {
		return []byte(cfg.RobotToMergePR.Token)
	})

	tmpl, err := newtemplateImpl(&cfg.Template)
	if err != nil {
		return nil, err
	}

	return &pullRequestImpl{
		cli:          cli,
		cfg:          *cfg,
		template:     tmpl,
		cliToMergePR: robot,
	}, nil
}

type iClient interface {
	CreatePullRequest(org, repo, title, body, head, base string, canModify bool) (sdk.PullRequest, error)
	GetGiteePullRequest(org, repo string, number int32) (sdk.PullRequest, error)
	ClosePR(org, repo string, number int32) error
	CreatePRComment(org, repo string, number int32, comment string) error
}

type clientToMergePR interface {
	MergePR(owner, repo string, number int32, opt sdk.PullRequestMergePutParam) error
}

type pullRequestImpl struct {
	cli          iClient
	cfg          Config
	template     templateImpl
	cliToMergePR clientToMergePR
}

func (impl *pullRequestImpl) Create(pkg *domain.SoftwarePkg) (domain.PullRequest, error) {
	if err := impl.createBranch(pkg); err != nil {
		return domain.PullRequest{}, err
	}

	return impl.createPR(pkg)
}

func (impl *pullRequestImpl) Merge(prNum int) error {
	org := impl.cfg.PR.Org
	repo := impl.cfg.PR.Repo

	v, err := impl.cli.GetGiteePullRequest(org, repo, int32(prNum))
	if err != nil {
		return err
	}

	if v.State != sdk.StatusOpen {
		return nil
	}

	return impl.cliToMergePR.MergePR(
		org, repo, int32(prNum), sdk.PullRequestMergePutParam{},
	)
}

func (impl *pullRequestImpl) Close(prNum int) error {
	org := impl.cfg.PR.Org
	repo := impl.cfg.PR.Repo

	prDetail, err := impl.cli.GetGiteePullRequest(org, repo, int32(prNum))
	if err != nil {
		return err
	}

	if prDetail.State != sdk.StatusOpen {
		return nil
	}

	return impl.cli.ClosePR(org, repo, int32(prNum))
}

func (impl *pullRequestImpl) Comment(prNum int, content string) error {
	return impl.cli.CreatePRComment(
		impl.cfg.PR.Org, impl.cfg.PR.Repo,
		int32(prNum), content,
	)
}
