package app

import (
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain/email"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain/message"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain/repository"
)

type PullRequestService interface {
	HandleCI(cmd *CmdToHandleCI) error
	HandleRepoCreated(*domain.PullRequest, string) error
	HandlePRMerged(cmd *CmdToHandlePRMerged) error
	HandlePRClosed(cmd *CmdToHandlePRClosed) error
}

func NewPullRequestService(
	r repository.PullRequest, p message.SoftwarePkgMessage, e email.Email,
) *pullRequestService {
	return &pullRequestService{
		repo:     r,
		producer: p,
		email:    e,
	}
}

type pullRequestService struct {
	repo     repository.PullRequest
	producer message.SoftwarePkgMessage
	email    email.Email
}

func (s *pullRequestService) HandleCI(cmd *CmdToHandleCI) error {
	pr, err := s.repo.Find(cmd.PRNum)
	if err != nil {
		return err
	}

	if !cmd.isSuccess() {
		if err = s.email.Send(pr.Link); err != nil {
			return err
		}
	}

	e := domain.NewPRCIFinishedEvent(&pr, cmd.FailedReason)
	return s.producer.NotifyCIResult(&e)
}

func (s *pullRequestService) HandleRepoCreated(pr *domain.PullRequest, url string) error {
	e := domain.NewRepoCreatedEvent(pr, url)
	if err := s.producer.NotifyRepoCreatedResult(&e); err != nil {
		return err
	}

	return s.repo.Remove(pr.Num)
}

func (s *pullRequestService) HandlePRMerged(cmd *CmdToHandlePRMerged) error {
	pr, err := s.repo.Find(cmd.PRNum)
	if err != nil {
		return err
	}

	e := domain.NewPRMergedEvent(&pr, cmd.ApprovedBy)
	if err = s.producer.NotifyPRMerged(&e); err != nil {
		return err
	}

	pr.SetMerged()

	return s.repo.Save(&pr)
}

func (s *pullRequestService) HandlePRClosed(cmd *CmdToHandlePRClosed) error {
	pr, err := s.repo.Find(cmd.PRNum)
	if err != nil {
		return err
	}

	e := domain.NewPRClosedEvent(&pr, cmd.Reason, cmd.RejectedBy)
	if err = s.producer.NotifyPRClosed(&e); err != nil {
		return err
	}

	return s.repo.Remove(pr.Num)
}
