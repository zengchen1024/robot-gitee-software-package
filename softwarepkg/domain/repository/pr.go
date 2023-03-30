package repository

import "github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain"

type PullRequest interface {
	Add(*domain.PullRequest) error
	Save(*domain.PullRequest) error
	Find(int) (domain.PullRequest, error)
	FindAll() ([]domain.PullRequest, error)
	Remove(int) error
}
