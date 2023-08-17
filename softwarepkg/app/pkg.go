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

type PackageService interface {
	HandleCreatePR(cmd *CmdToHandleNewPkg) error
	HandleCI(cmd *CmdToHandleCI) error
	HandleRepoCreated(*domain.SoftwarePkg, string) error
	HandlePRMerged(cmd *CmdToHandlePRMerged) error
	HandlePRClosed(cmd *CmdToHandlePRClosed) error
	HandlePushCode(*domain.SoftwarePkg) error
}

func NewPackageService(
	r repository.SoftwarePkg,
	p message.SoftwarePkgMessage,
	e email.Email,
	c pullrequest.PullRequest,
	cd code.Code,
) *packageService {
	return &packageService{
		repo:     r,
		producer: p,
		email:    e,
		prCli:    c,
		code:     cd,
	}
}

type packageService struct {
	repo     repository.SoftwarePkg
	producer message.SoftwarePkgMessage
	email    email.Email
	prCli    pullrequest.PullRequest
	code     code.Code
}

func (s *packageService) HandleCreatePR(cmd *CmdToHandleNewPkg) error {
	pr, err := s.prCli.Create(cmd)
	if err != nil {
		return err
	}

	cmd.PullRequest = pr
	cmd.SetPkgStatusPRCreated()

	return s.repo.Save(cmd)
}

func (s *packageService) HandleCI(cmd *CmdToHandleCI) error {
	pkg, err := s.repo.Find(cmd.PRNum)
	if err != nil {
		return err
	}

	if cmd.isSuccess() {
		if err := s.mergePR(pkg); err != nil {
			return s.notifyException(&pkg, err.Error())
		}
	} else {
		if !cmd.isPkgExisted() {
			return s.notifyException(&pkg, cmd.FailedReason)
		}
		s.closePR(pkg)
	}

	e := domain.PRCIFinishedEvent{
		PkgId:        pkg.Id,
		RelevantPR:   pkg.PullRequest.Link,
		RepoLink:     cmd.RepoLink,
		FailedReason: cmd.FailedReason,
	}

	return s.producer.NotifyCIResult(&e)
}

func (s *packageService) mergePR(pkg domain.SoftwarePkg) error {
	if err := s.prCli.Merge(pkg.PullRequest.Num); err != nil {
		return fmt.Errorf("merge pr(%d) failed: %s", pkg.PullRequest.Num, err.Error())
	}

	pkg.SetPkgStatusPRMerged()

	if err := s.repo.Save(&pkg); err != nil {
		logrus.Errorf("save pr(%d) failed: %s", pkg.PullRequest.Num, err.Error())
	}

	return nil
}

func (s *packageService) closePR(pkg domain.SoftwarePkg) {
	if err := s.prCli.Close(pkg.PullRequest.Num); err != nil {
		logrus.Errorf("close pr/%d failed: %s", pkg.PullRequest.Num, err.Error())
	}

	if err := s.repo.Remove(pkg.PullRequest.Num); err != nil {
		logrus.Errorf("remove pr/%d failed: %s", pkg.PullRequest.Num, err.Error())
	}
}

func (s *packageService) HandleRepoCreated(pkg *domain.SoftwarePkg, url string) error {
	e := domain.RepoCreatedEvent{
		PkgId:    pkg.Id,
		Platform: domain.PlatformGitee,
		RepoLink: url,
	}
	if err := s.producer.NotifyRepoCreatedResult(&e); err != nil {
		return err
	}

	pkg.SetPkgStatusRepoCreated()

	return s.repo.Save(pkg)
}

func (s *packageService) HandlePushCode(pkg *domain.SoftwarePkg) error {
	repoUrl, err := s.code.Push(pkg)
	if err != nil {
		logrus.Errorf("pkgId %s push code err: %s", pkg.Id, err.Error())

		return err
	}

	e := domain.CodePushedEvent{
		PkgId:    pkg.Id,
		Platform: domain.PlatformGitee,
		RepoLink: repoUrl,
	}

	if err = s.producer.NotifyCodePushedResult(&e); err != nil {
		return err
	}

	return s.repo.Remove(pkg.PullRequest.Num)
}

func (s *packageService) HandlePRMerged(cmd *CmdToHandlePRMerged) error {
	pkg, err := s.repo.Find(cmd.PRNum)
	if err != nil {
		return err
	}

	if pkg.IsPkgStatusMerged() {
		return nil
	}

	e := domain.PRCIFinishedEvent{
		PkgId:      pkg.Id,
		RelevantPR: pkg.PullRequest.Link,
	}

	if err = s.producer.NotifyCIResult(&e); err != nil {
		return err
	}

	pkg.SetPkgStatusPRMerged()

	return s.repo.Save(&pkg)
}

func (s *packageService) HandlePRClosed(cmd *CmdToHandlePRClosed) error {
	pkg, err := s.repo.Find(cmd.PRNum)
	if err != nil {
		return err
	}

	subject := fmt.Sprintf(
		"the pr of software package was closed by: %s",
		cmd.RejectedBy,
	)
	content := s.emailContent(pkg.PullRequest.Link)

	if err = s.email.Send(subject, content); err != nil {
		return fmt.Errorf("send email failed: %s", err.Error())
	}

	pkg.SetPkgStatusException()

	return s.repo.Save(&pkg)
}

func (s *packageService) emailContent(url string) string {
	return fmt.Sprintf("th pr url is: %s", url)
}

func (s *packageService) notifyException(
	pkg *domain.SoftwarePkg, reason string,
) error {
	subject := fmt.Sprintf(
		"the ci of software package check failed: %s",
		reason,
	)
	content := s.emailContent(pkg.PullRequest.Link)

	if err := s.email.Send(subject, content); err != nil {
		return fmt.Errorf("send email failed: %s", err.Error())
	}

	pkg.SetPkgStatusException()

	if err := s.repo.Save(pkg); err != nil {
		return fmt.Errorf("save pkg when exception error: %s", err.Error())
	}

	return nil
}
