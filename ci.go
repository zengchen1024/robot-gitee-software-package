package main

import (
	sdk "github.com/opensourceways/go-gitee/gitee"

	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/app"
)

func (bot *robot) handleCILabel(e *sdk.PullRequestEvent, cfg *botConfig) error {
	labels := e.PullRequest.LabelsToSet()

	if labels.Has(cfg.CILabel.Success) {
		cmd := bot.ciCmd(e.Number, "")
		if err := bot.prService.HandleCI(cmd); err != nil {
			return err
		}
	}

	if labels.Has(cfg.CILabel.Fail) {
		cmd := bot.ciCmd(e.Number, "ci check failed")
		if err := bot.prService.HandleCI(cmd); err != nil {
			return err
		}
	}

	return nil
}

func (bot *robot) ciCmd(num int64, reason string) *app.CmdToHandleCI {
	return &app.CmdToHandleCI{
		PRNum:        int(num),
		FailedReason: reason,
	}
}
