package app

import (
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain/pullrequest"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain/repository"
)

type MessageService interface {
	NewPkg(*CmdToHandleNewPkg) error
}

func NewMessageService(repo repository.SoftwarePkg, prCli pullrequest.PullRequest,
) *messageService {
	return &messageService{
		repo:  repo,
		prCli: prCli,
	}
}

type messageService struct {
	repo  repository.SoftwarePkg
	prCli pullrequest.PullRequest
}

func (s *messageService) NewPkg(cmd *CmdToHandleNewPkg) error {
	cmd.SetPkgStatusInitialized()

	return s.repo.Add(cmd)
}
