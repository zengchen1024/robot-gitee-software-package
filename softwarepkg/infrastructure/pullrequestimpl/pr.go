package pullrequestimpl

import (
	sdk "github.com/opensourceways/go-gitee/gitee"

	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain"
)

func NewPullRequestImpl(cli iClient, cfg Config) (impl *pullRequestImpl, err error) {
	v, err := cli.GetBot()
	if err != nil {
		return
	}

	impl = &pullRequestImpl{
		cli:        cli,
		cfg:        cfg,
		robotLogin: v.Login,
	}

	return
}

type pullRequestImpl struct {
	cli        iClient
	cfg        Config
	robotLogin string
}

type iClient interface {
	GetBot() (sdk.User, error)
	CreatePullRequest(org, repo, title, body, head, base string, canModify bool) (sdk.PullRequest, error)
	GetGiteePullRequest(org, repo string, number int32) (sdk.PullRequest, error)
	MergePR(owner, repo string, number int32, opt sdk.PullRequestMergePutParam) error
	ClosePR(org, repo string, number int32) error
	CreatePRComment(org, repo string, number int32, comment string) error
}

func (impl *pullRequestImpl) Create(pkg *domain.SoftwarePkg) (pr domain.PullRequest, err error) {
	if err = impl.initRepo(pkg); err != nil {
		return
	}

	if err = impl.newBranch(pkg); err != nil {
		return
	}

	if err = impl.modifyFiles(pkg); err != nil {
		return
	}

	if err = impl.commit(pkg); err != nil {
		return
	}

	v, err := impl.submit(pkg)
	if err == nil {
		pr.Num = int(v.Number)
		pr.Link = v.HtmlUrl
	}

	return
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

	return impl.cli.MergePR(org, repo, int32(prNum), sdk.PullRequestMergePutParam{})
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
