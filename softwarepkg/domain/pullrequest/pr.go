package pullrequest

import "github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain"

type PullRequest interface {
	Create(*domain.SoftwarePkg) (domain.PullRequest, error)
	Merge(*domain.PullRequest) error
	Close(*domain.PullRequest) error
}
