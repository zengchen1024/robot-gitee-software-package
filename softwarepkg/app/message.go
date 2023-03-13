package app

import (
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain/pullrequest"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain/repository"
)

type MessageService interface {
	CreatePR(cmd *CmdToCreatePR) error
	MergePR(cmd *CmdToMergePR) error
}

type messageService struct {
	repo  repository.PullRequest
	prCli pullrequest.PullRequest
}

func (s *messageService) CreatePR(cmd *CmdToCreatePR) error {
	pr, err := s.prCli.Create(cmd)
	if err != nil {
		return err
	}

	return s.repo.Add(&pr)
}

func (s *messageService) MergePR(cmd *CmdToMergePR) error {
	pr, err := s.repo.Find(cmd.PRNum)
	if err != nil {
		return err
	}

	if err = s.prCli.Merge(&pr); err != nil {
		return err
	}

	//TODO monitor create repo

	return nil
}
