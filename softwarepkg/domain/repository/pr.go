package repository

import "github.com/opensourceways/robot-gitee-software-package/pullrequest/domain"

type PullRequest interface {
	Add(*domain.PullRequest) (int, error)
	Find(int) (domain.PullRequest, error)
}
