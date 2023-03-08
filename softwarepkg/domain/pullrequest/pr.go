package pullrequest

import "github.com/opensourceways/robot-gitee-software-package/pullrequest/domain"

type PullRequest interface {
	Create(*domain.SoftwarePkg) (int, error)
	Merge(*domain.PullRequest) error
	Close(*domain.PullRequest) error
}
