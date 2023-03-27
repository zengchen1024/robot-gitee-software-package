package app

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain/email"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain/message"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain/pullrequest"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain/repository"
)

type PullRequestService interface {
	HandleCI(cmd *CmdToHandleCI) error
	HandleRepoCreated(*domain.PullRequest, string) error
	HandlePRMerged(cmd *CmdToHandlePRMerged) error
	HandlePRClosed(cmd *CmdToHandlePRClosed) error
}

func NewPullRequestService(
	r repository.PullRequest,
	p message.SoftwarePkgMessage,
	e email.Email,
	c pullrequest.PullRequest,
) *pullRequestService {
	return &pullRequestService{
		repo:     r,
		producer: p,
		email:    e,
		prCli:    c,
	}
}

type pullRequestService struct {
	repo     repository.PullRequest
	producer message.SoftwarePkgMessage
	email    email.Email
	prCli    pullrequest.PullRequest
}

func (s *pullRequestService) HandleCI(cmd *CmdToHandleCI) error {
	pr, err := s.repo.Find(cmd.PRNum)
	if err != nil {
		return err
	}

	if cmd.isSuccess() {
		s.handleSuccess(pr, cmd)
	} else {
		s.ciFailedSendEmail(&pr, cmd)
	}

	if cmd.isPkgExisted() {
		if err = s.prCli.Close(&pr); err != nil {
			logrus.Errorf("close pr failed: %s", err.Error())
		}
	}

	e := domain.NewPRCIFinishedEvent(&pr, cmd.FailedReason, cmd.RepoLink)
	return s.producer.NotifyCIResult(&e)
}

func (s *pullRequestService) handleSuccess(pr domain.PullRequest, cmd *CmdToHandleCI) {
	if err := s.prCli.Merge(&pr); err != nil {
		cmd.FailedReason = "merge pr failed"
		logrus.Errorf("merge pr failed: %s", err.Error())

		return
	}

	pr.SetMerged()

	if err := s.repo.Save(&pr); err != nil {
		cmd.FailedReason = "save pr failed"
		logrus.Errorf("save pr failed: %s", err.Error())
	}
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

func (s *pullRequestService) ciFailedSendEmail(
	pr *domain.PullRequest, cmd *CmdToHandleCI,
) {
	subject := fmt.Sprintf(
		"the ci of software package check failed: %s",
		cmd.FailedReason,
	)
	content := fmt.Sprintf("th pr url is: %s", pr.Link)

	if err := s.email.Send(subject, content); err != nil {
		logrus.Errorf("send email failed: %s", err.Error())
	}
}
