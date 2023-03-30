package main

import (
	sdk "github.com/opensourceways/go-gitee/gitee"

	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/app"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain/repository"
)

func (bot *robot) handleCILabel(e *sdk.PullRequestEvent, cfg *botConfig) error {
	pkg, err := bot.repo.Find(int(e.Number))
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

	if labels.Has(cfg.CILabel.Success) {
		return bot.prService.HandleCI(&cmd)
	}

	if labels.Has(cfg.CILabel.Fail) {
		cmd.FailedReason = "ci check failed"

		if v, err := bot.cli.GetRepo(bot.PkgSrcOrg, pkg.Name); err == nil {
			cmd.RepoLink = v.HtmlUrl
			cmd.FailedReason = "package already exists"
		}

		return bot.prService.HandleCI(&cmd)
	}

	return nil
}
