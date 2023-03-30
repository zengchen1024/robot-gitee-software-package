package app

import (
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain/pullrequest"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain/repository"
)

type MessageService interface {
	NewPkg(*CmdToHandleNewPkg) error
}

func NewMessageService(repo repository.PullRequest, prCli pullrequest.PullRequest,
) *messageService {
	return &messageService{
		repo:  repo,
		prCli: prCli,
	}
}

type messageService struct {
	repo  repository.PullRequest
	prCli pullrequest.PullRequest
}

func (s *messageService) NewPkg(cmd *CmdToHandleNewPkg) error {
	pr, err := s.prCli.Create(cmd)
	if err != nil {
		return err
	}

	pr.SetStatusInitialized()

	return s.repo.Add(&pr)
}
