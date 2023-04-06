package community

import (
	sdk "github.com/opensourceways/go-gitee/gitee"

	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/app"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain/repository"
)

func (impl *eventHandler) handleCILabel(e *sdk.PullRequestEvent) error {
	pkg, err := impl.repo.Find(int(e.Number))
	if err != nil {
		if repository.IsErrorResourceNotFound(err) {
			err = nil
		}

		return err
	}

	cmd := app.CmdToHandleCI{
		PRNum: int(e.Number),
	}

	labels := e.PullRequest.LabelsToSet()
	cfg := &impl.cfg

	if labels.Has(cfg.CISuccessLabel) {
		return impl.service.HandleCI(&cmd)
	}

	if labels.Has(cfg.CIFailureLabel) {
		cmd.FailedReason = "ci check failed"

		if v, err := impl.cli.GetRepo(cfg.PkgOrg, pkg.Name); err == nil {
			cmd.RepoLink = v.HtmlUrl
			cmd.FailedReason = "package already exists"
		}

		return impl.service.HandleCI(&cmd)
	}

	return nil
}
