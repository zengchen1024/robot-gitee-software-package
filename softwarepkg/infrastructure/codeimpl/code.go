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

	repoUrl := fmt.Sprintf(
		"https://gitee.com/%s/",
		cfg.PkgSrcOrg,
	)

	return &codeImpl{
		gitUrl:  gitUrl,
		repoUrl: repoUrl,
		script:  cfg.ShellScript,
	}
}

type codeImpl struct {
	gitUrl  string
	repoUrl string
	script  string
}

func (impl *codeImpl) Push(pkg *domain.SoftwarePkg) (string, error) {
	repoUrl := fmt.Sprintf("%s%s.git", impl.gitUrl, pkg.Name)

	params := []string{
		impl.script,
		repoUrl,
		pkg.Name,
		pkg.Importer.Name,
		pkg.Importer.Email,
		pkg.Application.SourceCode.SpecURL,
		pkg.Application.SourceCode.SrcRPMURL,
	}

	_, err, _ := utils.RunCmd(params...)
	if err != nil {
		logrus.Errorf(
			"run push code shell, err=%s, params=%v",
			err.Error(), params[:len(params)-1],
		)
	}

	return impl.repoUrl + pkg.Name, err
}
