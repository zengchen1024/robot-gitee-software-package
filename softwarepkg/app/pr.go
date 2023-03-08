package app

import (
	"github.com/opensourceways/robot-gitee-software-package/pullrequest/domain"
	"github.com/opensourceways/robot-gitee-software-package/pullrequest/domain/message"
	"github.com/opensourceways/robot-gitee-software-package/pullrequest/domain/pullrequest"
	"github.com/opensourceways/robot-gitee-software-package/pullrequest/domain/repository"
)

type PullRequestSerivce interface{}

type pullRequestSerivce struct {
	repo     repository.PullRequest
	prCli    pullrequest.PullRequest
	producer message.SoftwarePkgMessage
}

func (s *pullRequestSerivce) CreatePR() error {
	// create pr
	// save the pr

	return nil
}

func (s *pullRequestSerivce) HandleCI(num int, success bool, failedReason string) error {
	pr, err := s.repo.Find(num)
	if err != nil {
		return err
	}

	e := domain.NewPRCIFinishedEvent(&pr)

	s.producer.NotifyCIResult(&e)
	return nil
}
