package pullrequestimpl

import (
	sdk "github.com/opensourceways/go-gitee/gitee"

	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain"
)

type pullRequestImpl struct {
	cli        iClient
	cfg        Config
	pkg        *domain.SoftwarePkg
	robotLogin string
}

type iClient interface {
	GetBot() (sdk.User, error)
	CreatePullRequest(org, repo, title, body, head, base string, canModify bool) (sdk.PullRequest, error)
	GetGiteePullRequest(org, repo string, number int32) (sdk.PullRequest, error)
	MergePR(owner, repo string, number int32, opt sdk.PullRequestMergePutParam) error
}

func (impl *pullRequestImpl) Create(pkg *domain.SoftwarePkg) (pr domain.PullRequest, err error) {
	impl.pkg = pkg

	if err = impl.initRepo(); err != nil {
		return
	}

	if err = impl.newBranch(); err != nil {
		return
	}

	if err = impl.modifyFiles(); err != nil {
		return
	}

	if err = impl.commit(); err != nil {
		return
	}

	return impl.submit()
}

func (impl *pullRequestImpl) Merge(pr *domain.PullRequest) error {
	org := impl.cfg.PR.Org
	repo := impl.cfg.PR.Repo

	v, err := impl.cli.GetGiteePullRequest(org, repo, int32(pr.Num))
	if err != nil {
		return err
	}

	if v.State != sdk.StatusOpen {
		return nil
	}

	return impl.cli.MergePR(org, repo, int32(pr.Num), sdk.PullRequestMergePutParam{})
}

func (impl *pullRequestImpl) Close(*domain.PullRequest) error {
	return nil
}
