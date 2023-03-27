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

	labels := e.PullRequest.LabelsToSet()

	if labels.Has(cfg.CILabel.Success) {
		cmd := bot.ciCmd(e.Number, "", "")
		if err := bot.prService.HandleCI(&cmd); err != nil {
			return err
		}
	}

	if labels.Has(cfg.CILabel.Fail) {
		cmd := bot.ciCmd(e.Number, "", "ci check failed")
		if v, err := bot.cli.GetRepo(bot.PkgSrcOrg, dpr.Pkg.Name); err == nil {
			cmd.RepoLink = v.HtmlUrl
			cmd.FailedReason = "package already exists"
		}

		if err = bot.prService.HandleCI(&cmd); err != nil {
			return err
		}
	}

	return nil
}

func (bot *robot) ciCmd(num int64, link, reason string) app.CmdToHandleCI {
	return app.CmdToHandleCI{
		PRNum:        int(num),
		RepoLink:     link,
		FailedReason: reason,
	}
}
