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
		if err := s.mergePR(pr); err != nil {
			cmd.FailedReason = err.Error()
			s.notifyException(&pr, cmd)
		}
	} else {
		if cmd.isPkgExisted() {
			if err = s.prCli.Close(&pr); err != nil {
				logrus.Errorf("close pr failed: %s", err.Error())
			}
		} else {
			s.notifyException(&pr, cmd)
		}
	}

	e := domain.NewPRCIFinishedEvent(&pr, cmd.FailedReason, cmd.RepoLink)

	return s.producer.NotifyCIResult(&e)
}

func (s *pullRequestService) mergePR(pr domain.PullRequest) error {
	if err := s.prCli.Merge(&pr); err != nil {
		return fmt.Errorf("merge pr(%d) failed: %s", pr.Num, err.Error())
	}

	pr.SetMerged()

	if err := s.repo.Save(&pr); err != nil {
		logrus.Errorf("save pr(%d) failed: %s", pr.Num, err.Error())
	}

	return nil
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

func (s *pullRequestService) notifyException(
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
