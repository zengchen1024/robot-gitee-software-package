package main

import (
	"strings"

	sdk "github.com/opensourceways/go-gitee/gitee"

	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/app"
)

func (bot *robot) handleCILabel(e *sdk.PullRequestEvent, cfg *botConfig) error {
	dpr, err := bot.repo.Find(int(e.Number))
	if err != nil {
		if strings.Contains(err.Error(), "row not found") {
			return nil
		} else {
			return err
		}
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

		if v, err := bot.cli.GetRepo(bot.PkgSrcOrg, dpr.Pkg.Name); err == nil {
			cmd.RepoLink = v.HtmlUrl
			cmd.FailedReason = "package already exists"
		}

		return bot.prService.HandleCI(&cmd)
	}

	return nil
}
