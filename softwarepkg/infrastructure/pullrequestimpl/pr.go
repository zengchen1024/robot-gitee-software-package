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

func (impl *pullRequestImpl) Merge(*domain.PullRequest) error {
	return nil
}

func (impl *pullRequestImpl) Close(*domain.PullRequest) error {
	return nil
}
