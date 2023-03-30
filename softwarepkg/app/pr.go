package app

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain/code"
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
	HandlePushCode(pr *domain.PullRequest) error
}

func NewPullRequestService(
	r repository.PullRequest,
	p message.SoftwarePkgMessage,
	e email.Email,
	c pullrequest.PullRequest,
	cd code.Code,
) *pullRequestService {
	return &pullRequestService{
		repo:     r,
		producer: p,
		email:    e,
		prCli:    c,
		code:     cd,
	}
}

type pullRequestService struct {
	repo     repository.PullRequest
	producer message.SoftwarePkgMessage
	email    email.Email
	prCli    pullrequest.PullRequest
	code     code.Code
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
			s.closePR(pr)
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

	pr.SetStatusMerged()

	if err := s.repo.Save(&pr); err != nil {
		logrus.Errorf("save pr(%d) failed: %s", pr.Num, err.Error())
	}

	return nil
}

func (s *pullRequestService) closePR(pr domain.PullRequest) {
	if err := s.prCli.Close(&pr); err != nil {
		logrus.Errorf("close pr/%d failed: %s", pr.Num, err.Error())
	}

	if err := s.repo.Remove(pr.Num); err != nil {
		logrus.Errorf("remove pr/%d failed: %s", pr.Num, err.Error())
	}
}

func (s *pullRequestService) HandleRepoCreated(pr *domain.PullRequest, url string) error {
	pr.SetStatusRepoCreated()

	if err := s.repo.Save(pr); err != nil {
		return err
	}

	e := domain.NewRepoCreatedEvent(pr, url, "")

	return s.producer.NotifyRepoCreatedResult(&e)
}

func (s *pullRequestService) HandlePushCode(pr *domain.PullRequest) error {
	if err := s.code.Push(pr); err != nil {
		logrus.Errorf("pkgId %s push code err: %s", pr.Pkg.Id, err.Error())

		return err
	}

	e := domain.CodePushedEvent{
		PkgId:    pr.Pkg.Id,
		Platform: domain.PlatformGitee,
	}

	if err := s.producer.NotifyCodePushedResult(&e); err != nil {
		return err
	}

	return s.repo.Remove(pr.Num)
}

func (s *pullRequestService) HandlePRMerged(cmd *CmdToHandlePRMerged) error {
	pr, err := s.repo.Find(cmd.PRNum)
	if err != nil {
		return err
	}

	if pr.IsStatusMerged() {
		return nil
	}

	e := domain.PRCIFinishedEvent{
		PkgId:      pr.Pkg.Id,
		RelevantPR: pr.Link,
	}

	if err = s.producer.NotifyCIResult(&e); err != nil {
		return err
	}

	pr.SetStatusMerged()

	return s.repo.Save(&pr)
}

func (s *pullRequestService) HandlePRClosed(cmd *CmdToHandlePRClosed) error {
	pr, err := s.repo.Find(cmd.PRNum)
	if err != nil {
		return err
	}

	subject := fmt.Sprintf(
		"the pr of software package was closed by: %s",
		cmd.RejectedBy,
	)
	content := s.emailContent(pr.Link)

	if err = s.email.Send(subject, content); err != nil {
		logrus.Errorf("send email failed: %s", err.Error())
	}

	return nil
}

func (s *pullRequestService) emailContent(url string) string {
	return fmt.Sprintf("th pr url is: %s", url)
}

func (s *pullRequestService) notifyException(
	pr *domain.PullRequest, cmd *CmdToHandleCI,
) {
	subject := fmt.Sprintf(
		"the ci of software package check failed: %s",
		cmd.FailedReason,
	)
	content := s.emailContent(pr.Link)

	if err := s.email.Send(subject, content); err != nil {
		logrus.Errorf("send email failed: %s", err.Error())
	}
}
