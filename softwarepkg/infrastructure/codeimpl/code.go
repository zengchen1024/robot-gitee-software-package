package codeimpl

import (
	"fmt"

	"github.com/opensourceways/server-common-lib/utils"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain"
)

func NewCodeImpl(cfg Config) *codeImpl {
	gitUrl := fmt.Sprintf(
		"https://%s:%s@gitee.com/%s/",
		cfg.Robot.Username,
		cfg.Robot.Token,
		cfg.PkgSrcOrg,
	)

	return &codeImpl{
		gitUrl: gitUrl,
		script: cfg.ShellScript,
	}
}

type codeImpl struct {
	gitUrl string
	script string
}

func (impl *codeImpl) Push(pr *domain.PullRequest) error {
	repoUrl := fmt.Sprintf("%s%s.git", impl.gitUrl, pr.Pkg.Name)

	params := []string{
		impl.script,
		repoUrl,
		pr.Pkg.Name,
		pr.ImporterName,
		pr.ImporterEmail,
		pr.SrcCode.SpecURL,
		pr.SrcCode.SrcRPMURL,
	}

	_, err, _ := utils.RunCmd(params...)
	if err != nil {
		logrus.Errorf(
			"run push code shell, err=%s, params=%v",
			err.Error(), params[:len(params)-1],
		)
	}

	return err
}
